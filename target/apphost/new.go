package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	sig2 "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/target"
	"time"
)

type Factory struct {
	invoke target.Dispatch
}

func NewFactory(invoke target.Dispatch) *Factory {
	return &Factory{invoke: invoke}
}

func newAdapter(ctx context.Context, pkg string) *Adapter {
	if pkg == "" {
		panic("package is empty")
	}
	a := &Adapter{}
	a.port = port.New(pkg)
	a.log = plog.Get(ctx).Type(a).Set(&ctx)

	a.listeners = mem.NewCache[*Listener]()
	a.connections = mem.NewCache[*Conn]()

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

func (f Factory) NewAdapter(ctx context.Context, portal target.Portal) target.Apphost {
	flat := newAdapter(ctx, portal.Manifest().Package)
	go func() {
		sleep := time.Duration(3) * time.Millisecond
		time.Sleep(sleep)
		flat.events.Push(target.ApphostEvent{Type: target.ApphostDisconnect})
	}()
	return NewInvoker(ctx, flat, f.invoke)
}

var ConnectionsThreshold = -1

func (f Factory) WithTimeout(ctx context.Context, portal target.Portal) target.Apphost {
	manifest := portal.Manifest()
	flat := newAdapter(ctx, manifest.Package)

	if manifest.Env.Timeout > -1 && ConnectionsThreshold >= 0 {
		go func() {
			duration := 5 * time.Second
			if manifest.Env.Timeout > 0 {
				duration = time.Duration(manifest.Env.Timeout) * time.Millisecond
			}
			timeout := NewTimout(duration, func() {
				_ = sig2.Interrupt()
			})
			timeout.Enable(true)
			for range flat.Events().Subscribe(ctx) {
				activeConnections := flat.connections.Size()
				timeout.Enable(activeConnections <= ConnectionsThreshold)
			}
		}()
	}
	return NewInvoker(ctx, flat, f.invoke)
}
