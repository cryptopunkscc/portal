package apphost

import "C"
import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/registry"
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
		prefix: prefix,
	}
	if pkg != "" {
		a.pkg = []string{pkg}
	}
	a.log = plog.Get(ctx).Type(a).Set(&ctx)

	a.listeners = registry.New[*Listener]()
	a.connections = registry.New[*Conn]()

	a.listeners.OnChange(eventEmitter[*Listener](a.Events()))
	a.connections.OnChange(eventEmitter[*Conn](a.Events()))

	return a
}

func eventEmitter[T any](queue *sig.Queue[target.ApphostEvent]) func(ref string, conn T, added bool) {
	return func(ref string, conn T, added bool) {
		event := target.ApphostEvent{Ref: ref}
		switch v := any(conn).(type) {
		case *Conn:
			event.Port = v.conn.Query()
			event.Type = target.ApphostDisconnect
			if added {
				event.Type = target.ApphostConnect
			}
		case *Listener:
			event.Port = v.port
			event.Type = target.ApphostUnregister
			if added {
				event.Type = target.ApphostRegister
			}
		default:
			return
		}
		queue.Push(event)
	}
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

	go func() {
		for range flat.Events().Subscribe(ctx) {
			timeout.Enable(flat.connections.Size() == 0)
		}
	}()

	return NewInvoker(ctx, flat, f.invoke)
}
