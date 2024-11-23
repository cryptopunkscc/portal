package apphost

import "github.com/cryptopunkscc/portal/api/apphost"

func Rpc(client apphost.Client) RpcBase {
	return RpcBase{client: client}
}

type RpcBase struct {
	client apphost.Client
}
