package apphost

import (
	apphost "github.com/cryptopunkscc/portal/core/apphost/rpc"
)

func (a *Adapter) Rpc() apphost.Rpc {
	return apphost.Rpc{Apphost: a}
}
