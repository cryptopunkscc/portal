package apphost

import (
	apphost "github.com/cryptopunkscc/portal/core/apphost/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func (a *Adapter) Rpc() *apphost.Rpc {
	return &apphost.Rpc{
		Apphost: a,
		Log:     a.Log,
	}
}

var _ rpc.Rpc = &Adapter{}

func (a *Adapter) Format(name string) rpc.Rpc {
	return a.Rpc().Format(name)
}

func (a *Adapter) Conn(target, query string) (rpc.Conn, error) {
	return a.Rpc().Conn(target, query)
}

func (a *Adapter) Request(target string, query ...string) rpc.Conn {
	return a.Rpc().Request(target, query...)
}
