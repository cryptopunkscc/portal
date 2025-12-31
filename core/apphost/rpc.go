package apphost

import (
	apphost "github.com/cryptopunkscc/portal/core/apphost/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func (a *Adapter) Rpc() rpc.Rpc {
	return &apphost.Rpc{
		Apphost: a,
		Log:     a.Log,
	}
}

var _ rpc.Rpc = &Adapter{}

func (a *Adapter) Router(handler cmd.Handler) rpc.Router {
	return a.Rpc().Router(handler)
}
