package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"log"
	"syscall"
	"time"
)

func NewAdapter(ctx context.Context, serve Invoke, prefix ...string) *Invoker {
	log.Println("NewAdapter", prefix)
	flat := &Adapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]*Conn{},
		prefix:      append([]string{}, prefix...),
	}
	return NewInvoker(ctx, flat, serve)
}

func WithTimeout(ctx context.Context, serve Invoke, prefix ...string) *Invoker {
	log.Println("NewAdapterWithTimeout", prefix)
	timeout := NewTimout(5*time.Second, func() {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})
	flat := &Adapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]*Conn{},
		prefix:      append([]string{}, prefix...),
		onIdle:      timeout.Enable,
	}
	return NewInvoker(ctx, flat, serve)
}
