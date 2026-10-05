package core

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/sagernet/sing-box/option"
)

// PolicyUpdate is prepared while the previous policy serves traffic. The
// caller serializes lifecycle changes, commits storage, then publishes. Abort
// releases every prepared resource if validation or the database commit fails.
type PolicyUpdate struct {
	box     *Box
	runtime *policyRuntime
	options option.Options
	impact  policyImpact
	once    sync.Once
}

func (u *PolicyUpdate) Abort() {
	u.once.Do(func() {
		if u.runtime != nil {
			u.runtime.close()
		}
	})
}
func (u *PolicyUpdate) Commit() {
	u.once.Do(func() {
		if u.runtime != nil {
			u.box.policy.publish(u.runtime.generation, u.impact)
		}
		u.box.applied = u.options
	})
}

func (c *Core) PreparePolicy(raw []byte) (*PolicyUpdate, error) {
	b, err := c.running()
	if err != nil {
		return nil, err
	}
	var options option.Options
	if err := options.UnmarshalJSONContext(c.ctx, raw); err != nil {
		return nil, err
	}
	if !sameOptions(policyBase(b.applied), policyBase(options)) {
		return nil, fmt.Errorf("日志、NTP、证书或实验功能变更需要停止核心后保存；当前拨号会话未中断")
	}
	update := &PolicyUpdate{box: b, options: options}
	if sameOptions(b.applied.Route, options.Route) && sameOptions(b.applied.DNS, options.DNS) && sameOptions(b.applied.Outbounds, options.Outbounds) && sameOptions(b.applied.HTTPClients, options.HTTPClients) {
		return update, nil
	}
	// These services hold native manager references outside the routing entry.
	// Refuse an incomplete reload instead of claiming it applied or redialling.
	if len(options.Endpoints) != 0 || len(options.Services) != 0 || len(options.CertificateProviders) != 0 || (options.Experimental != nil && (options.Experimental.ClashAPI != nil || options.Experimental.V2RayAPI != nil)) {
		return nil, fmt.Errorf("端点、核心服务或证书提供器启用时暂不支持策略热更新，请停止核心后保存")
	}
	if options.NTP != nil && options.NTP.Enabled {
		return nil, fmt.Errorf("核心 NTP 服务启用时暂不支持策略热更新，请停止核心后保存")
	}
	if options.Experimental != nil && options.Experimental.CacheFile != nil && options.Experimental.CacheFile.Enabled {
		return nil, fmt.Errorf("核心持久缓存启用时暂不支持策略热更新，请停止核心后保存")
	}
	for _, inbound := range b.inbound.Inbounds() {
		if inbound.Type() == "tun" {
			return nil, fmt.Errorf("TUN 入站的流卸载暂不支持策略热更新，请停止核心后保存")
		}
	}
	p, err := b.preparePolicy(options)
	if err != nil {
		return nil, fmt.Errorf("准备策略失败，原配置保持运行: %w", err)
	}
	update.runtime = p
	if !sameOptions(b.applied.Route, options.Route) || !sameOptions(b.applied.HTTPClients, options.HTTPClients) {
		// A route rule may depend on sniffed data, DNS answers or a remote
		// rule-set. Re-evaluating stored metadata would invent another router.
		update.impact.all = true
	} else {
		changed := changedOutboundTags(b.applied.Outbounds, options.Outbounds)
		b.policy.mu.Lock()
		current := b.policy.current.runtime
		b.policy.mu.Unlock()
		if current.outbound.Default().Tag() != p.outbound.Default().Tag() {
			changed[current.outbound.Default().Tag()] = true
		}
		// Closing a changed detour also invalidates selectors and other users of
		// it, even when their visible outbound tag itself did not change.
		for more := true; more; {
			more = false
			for _, manager := range []*policyRuntime{current, p} {
				for _, outbound := range manager.outbound.Outbounds() {
					if changed[outbound.Tag()] {
						continue
					}
					for _, dependency := range outbound.Dependencies() {
						if changed[dependency] {
							changed[outbound.Tag()] = true
							more = true
							break
						}
					}
				}
			}
		}
		update.impact.tags = changed
	}
	return update, nil
}

func policyBase(o option.Options) option.Options {
	o.RawMessage, o.CommentsSet = nil, nil
	o.DNS, o.Route, o.Outbounds, o.HTTPClients = nil, nil, nil, nil
	o.Inbounds, o.Endpoints, o.Services = nil, nil, nil
	return o
}

func sameOptions(a, b any) bool {
	// Protocol Options fields are json:"-" and use context-aware marshalers.
	// Compare decoded values, never encoding/json's incomplete representation.
	return reflect.DeepEqual(a, b)
}

func changedOutboundTags(before, after []option.Outbound) map[string]bool {
	previous := make(map[string]option.Outbound, len(before))
	changed := make(map[string]bool)
	for i, outbound := range before {
		tag := outbound.Tag
		if tag == "" {
			tag = fmt.Sprint(i)
		}
		previous[tag] = outbound
	}
	for i, outbound := range after {
		tag := outbound.Tag
		if tag == "" {
			tag = fmt.Sprint(i)
		}
		old, found := previous[tag]
		if !found || !sameOptions(old, outbound) {
			changed[tag] = true
		}
		delete(previous, tag)
	}
	for tag := range previous {
		changed[tag] = true
	}
	return changed
}
