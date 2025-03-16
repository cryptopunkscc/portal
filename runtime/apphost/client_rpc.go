package apphost

import (
	apphost "github.com/cryptopunkscc/portal/runtime/apphost/rpc"
)

func (a *Adapter) Rpc() apphost.Rpc {
	return apphost.Rpc{Apphost: a}
}
