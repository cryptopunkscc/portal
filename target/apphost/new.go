package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/go-astral-js/target"
	"syscall"
	"time"
)

type Factory struct {
	Invoke
	Prefix []string
}

func NewFactory(invoke Invoke, prefix ...string) *Factory {
	return &Factory{Invoke: invoke, Prefix: prefix}
}

func newAdapter(pkg string, prefix ...string) *Adapter {
	a := &Adapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]*Conn{},
		prefix:      prefix,
	}
	if pkg != "" {
		a.pkg = []string{pkg}
	}
	return a
}

func (f Factory) NewAdapter(ctx context.Context, pkg string) target.Apphost {
	flat := newAdapter(pkg, f.Prefix...)
	return NewInvoker(ctx, flat, f.Invoke)
}

func (f Factory) WithTimeout(ctx context.Context, pkg string) target.Apphost {
	timeout := NewTimout(5*time.Second, func() {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})
	flat := newAdapter(pkg, f.Prefix...)
	flat.onIdle = timeout.Enable
	return NewInvoker(ctx, flat, f.Invoke)
}
