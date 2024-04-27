package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"syscall"
	"time"
)

type Flat interface {
	Close() error
	Interrupt()
	Log(arg ...any)
	LogArr(arg []any)
	Sleep(duration int64)
	ServiceRegister(service string) (err error)
	ServiceClose(service string) (err error)
	ConnAccept(service string) (data string, err error)
	ConnClose(id string) (err error)
	ConnWrite(id string, data string) (err error)
	ConnRead(id string) (data string, err error)
	Query(identity string, query string) (data string, err error)
	QueryName(name string, query string) (data string, err error)
	Resolve(name string) (id string, err error)
	NodeInfo(identity string) (info NodeInfo, err error)
}

func NewAdapter(ctx context.Context) Flat {
	return &Invoker{
		ctx: ctx,
		Flat: &FlatAdapter{
			listeners:   map[string]*astral.Listener{},
			connections: map[string]*Conn{},
		},
	}
}

func WithTimeout(ctx context.Context) Flat {
	timeout := NewTimout(3*time.Second, func() {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})
	return &Invoker{
		ctx: ctx,
		Flat: &FlatAdapter{
			listeners:   map[string]*astral.Listener{},
			connections: map[string]*Conn{},
			onIdle:      timeout.Enable,
		},
	}
}
