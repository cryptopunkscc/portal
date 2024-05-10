package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"log"
	"syscall"
	"time"
)

type Flat interface {
	port(query string) string
	Prefix() []string
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

func NewAdapter(ctx context.Context, serve Invoke, prefix ...string) Flat {
	log.Println("NewAdapter", prefix)
	flat := &FlatAdapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]*Conn{},
		prefix:      append([]string{}, prefix...),
	}
	return NewInvoker(ctx, flat, serve)
}

func WithTimeout(ctx context.Context, serve Invoke, prefix ...string) Flat {
	log.Println("NewAdapterWithTimeout", prefix)
	timeout := NewTimout(5*time.Second, func() {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})
	flat := &FlatAdapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]*Conn{},
		prefix:      append([]string{}, prefix...),
		onIdle:      timeout.Enable,
	}
	return NewInvoker(ctx, flat, serve)
}
