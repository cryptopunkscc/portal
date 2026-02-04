package apphost

import (
	apphost "github.com/cryptopunkscc/portal/pkg/apphost/rpc"
	"github.com/cryptopunkscc/portal/pkg/util/rpc"
)

func (a *Adapter) Rpc() rpc.Rpc {
	return &apphost.Rpc{
		Log:      a.Log,
		Register: a.Register,
	}
}
