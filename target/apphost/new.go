package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"syscall"
	"time"
)

type Factory struct {
	invoke target.Dispatch
	prefix []string
}

func NewFactory(invoke target.Dispatch, prefix ...string) *Factory {
	return &Factory{invoke: invoke, prefix: prefix}
}

func newAdapter(ctx context.Context, pkg string, prefix ...string) *Adapter {
	a := &Adapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]*Conn{},
		prefix:      prefix,
	}
	if pkg != "" {
		a.pkg = []string{pkg}
	}
	a.log = plog.Get(ctx).Type(a).Set(&ctx)
	return a
}

func (f Factory) NewAdapter(ctx context.Context, pkg string) target.Apphost {
	flat := newAdapter(ctx, pkg, f.prefix...)
	return NewInvoker(ctx, flat, f.invoke)
}

func (f Factory) WithTimeout(ctx context.Context, pkg string) target.Apphost {
	timeout := NewTimout(5*time.Second, func() {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})
	flat := newAdapter(ctx, pkg, f.prefix...)
	flat.onIdle = timeout.Enable
	return NewInvoker(ctx, flat, f.invoke)
}
