package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/adapter/outbound"
	"github.com/sagernet/sing-box/common/httpclient"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/dns"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/protocol/direct"
	"github.com/sagernet/sing-box/route"
	R "github.com/sagernet/sing-box/route/rule"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/service"
)

// A policy owns no inbound listeners or VPN sessions. Its context and native
// managers form one generation: dialers must never mix DNS/outbound references
// from different configurations. Retired generations drain existing traffic.
type policyRuntime struct {
	generation *policyGeneration
	ctx        context.Context
	router     *route.Router
	dns        *dns.Router
	transports *dns.TransportManager
	outbound   *outbound.Manager
	network    *route.NetworkManager
	connection *route.ConnectionManager
	http       adapter.LifecycleService
	closeOnce  sync.Once
	closed     atomic.Bool
	closeErr   error
	logger     log.ContextLogger
	owned      []adapter.Lifecycle
}

func (p *policyRuntime) close() error {
	p.closeOnce.Do(func() {
		for i := len(p.owned) - 1; i >= 0; i-- {
			if err := p.owned[i].Close(); err != nil {
				p.closeErr = errors.Join(p.closeErr, fmt.Errorf("close %T: %w", p.owned[i], err))
			}
		}
		if p.closeErr != nil && p.logger != nil {
			p.logger.Error("retire policy: ", p.closeErr)
		}
		p.closed.Store(true)
	})
	return p.closeErr
}

func (b *Box) preparePolicy(options option.Options) (p *policyRuntime, err error) {
	ctx := service.ExtendContext(b.ctx)
	p = &policyRuntime{ctx: ctx, logger: b.logger}
	p.generation = &policyGeneration{runtime: p}
	defer func() {
		if err != nil {
			p.close()
		}
	}()
	r := common.PtrValueOrDefault(options.Route)
	d := common.PtrValueOrDefault(options.DNS)
	f := b.logFactory
	p.outbound = outbound.NewManager(f.NewLogger("outbound"), service.FromContext[adapter.OutboundRegistry](ctx), b.endpoint, r.Final)
	p.owned = append(p.owned, p.outbound)
	registry := service.FromContext[adapter.DNSTransportRegistry](ctx)
	p.transports = dns.NewTransportManager(f.NewLogger("dns/transport"), registry, p.outbound, d.Final)
	p.owned = append(p.owned, p.transports)
	service.MustRegister[adapter.OutboundManager](ctx, p.outbound)
	service.MustRegister[adapter.DNSTransportManager](ctx, p.transports)
	p.dns, err = dns.NewRouter(ctx, f, d)
	if err != nil {
		return p, err
	}
	p.owned = append(p.owned, p.dns)
	service.MustRegister[adapter.DNSRouter](ctx, &policyDNS{DNSRouter: p.dns, generation: p.generation, policy: b.policy})
	service.MustRegister[adapter.DNSRuleSetUpdateValidator](ctx, p.dns)
	p.network, err = route.NewNetworkManager(ctx, f.NewLogger("network"), r, d)
	if err != nil {
		return p, err
	}
	p.owned = append(p.owned, p.network)
	service.MustRegister[adapter.NetworkManager](ctx, p.network)
	p.connection = route.NewConnectionManager(f.NewLogger("connection"))
	p.owned = append(p.owned, p.connection)
	service.MustRegister[adapter.ConnectionManager](ctx, p.connection)
	h := httpclient.NewManager(ctx, f.NewLogger("httpclient"), options.HTTPClients, r.DefaultHTTPClient)
	p.http = h
	p.owned = append(p.owned, h)
	service.MustRegister[adapter.HTTPClientManager](ctx, h)
	p.router = route.NewRouter(ctx, f, r, d)
	p.owned = append(p.owned, p.router)
	service.MustRegister[adapter.Router](ctx, p.router)
	if err = p.router.Initialize(r.Rules, r.RuleSet); err != nil {
		return p, err
	}
	if err = p.dns.Initialize(d.Rules); err != nil {
		return p, err
	}
	for i, o := range d.Servers {
		tag := o.Tag
		if tag == "" {
			tag = fmt.Sprint(i)
		}
		if err = p.transports.Create(ctx, f.NewLogger("dns/"+tag), tag, o.Type, o.Options); err != nil {
			return p, err
		}
	}
	for i, o := range options.Outbounds {
		tag := o.Tag
		if tag == "" {
			tag = fmt.Sprint(i)
		}
		outCtx := adapter.WithContext(ctx, &adapter.InboundContext{Outbound: tag})
		if err = p.outbound.Create(outCtx, p.router, f.NewLogger("outbound/"+tag), tag, o.Type, o.Options); err != nil {
			return p, err
		}
	}
	p.outbound.Initialize(func() (adapter.Outbound, error) {
		return direct.NewOutbound(ctx, p.router, f.NewLogger("outbound/direct"), "direct", option.DirectOutboundOptions{})
	})
	p.transports.Initialize(func() (adapter.DNSTransport, error) {
		return registry.CreateDNSTransport(ctx, f.NewLogger("dns/local"), "local", C.DNSTypeLocal, &option.LocalDNSServerOptions{})
	})
	h.Initialize(func() (*httpclient.ManagedTransport, error) {
		return httpclient.NewTransport(ctx, f.NewLogger("httpclient"), "", option.HTTPClientOptions{DefaultOutbound: true})
	})
	p.router.AppendTracker(b.sessionTracker)
	// Stop background rule downloads/probes before disposing the DNS and
	// network services their dialers use. Until startup, construction failures
	// use the incremental reverse-construction cleanup above.
	p.owned = []adapter.Lifecycle{p.network, p.transports, p.dns, p.connection, p.outbound, h, p.router}
	for _, stage := range []adapter.StartStage{adapter.StartStateInitialize, adapter.StartStateStart, adapter.StartStatePostStart, adapter.StartStateStarted} {
		// Same dependency order as Box startup; no access service is restarted.
		if err = adapter.Start(ctx, b.logger, stage, p.outbound, p.transports, p.network, p.connection, h, p.router, p.dns); err != nil {
			return p, err
		}
	}
	for _, rule := range p.router.Rules() {
		var target string
		switch action := rule.Action().(type) {
		case *R.RuleActionRoute:
			target = action.Outbound
		case *R.RuleActionBypass:
			target = action.Outbound
		}
		if target != "" {
			if _, found := p.outbound.Outbound(target); !found {
				return p, fmt.Errorf("route references missing outbound %q", target)
			}
		}
	}
	if p.outbound.Default() == nil {
		return p, errors.New("missing default outbound")
	}
	return p, nil
}
