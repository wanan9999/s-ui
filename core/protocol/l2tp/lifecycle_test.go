//go:build with_gvisor

package l2tp

import (
	"context"
	"errors"
	"testing"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/log"
)

func TestWorkerExitHealthAndCleanup(t *testing.T) {
	for _, cause := range []error{errors.New("test read failure"), nil} {
		in := &Inbound{logger: log.NewNOPFactory().Logger()}
		cleaned := false
		in.runServer(func() error { return cause }, func() error { cleaned = true; return nil })
		if in.HealthError() == nil {
			t.Fatal("dead worker reports healthy")
		}
		if cause != nil && !errors.Is(in.HealthError(), cause) {
			t.Fatal("read failure was lost")
		}
		if !cleaned {
			t.Fatal("failed worker retains its sockets")
		}
		_ = in.Close()
	}
}

func TestAdministrativeCloseIsHealthy(t *testing.T) {
	in := &Inbound{logger: log.NewNOPFactory().Logger()}
	if err := in.Close(); err != nil {
		t.Fatal(err)
	}
	in.runServer(func() error { return nil }, func() error { t.Fatal("double-close on ordinary shutdown"); return nil })
	if in.HealthError() != nil {
		t.Fatal("administrative stop reported as a failure")
	}
}

func TestEmptyInboundRemainsHealthy(t *testing.T) {
	raw, err := newInbound(context.Background(), nil, log.NewNOPFactory().Logger(), "empty", Options{})
	if err != nil {
		t.Fatal(err)
	}
	in := raw.(*Inbound)
	if err := in.Start(adapter.StartStatePostStart); err != nil {
		t.Fatal(err)
	}
	if in.HealthError() != nil {
		t.Fatal(in.HealthError())
	}
	_ = in.Close()
}
