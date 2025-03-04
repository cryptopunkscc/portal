package apphost

import (
	"context"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

var ServeHandler = cmd.Handler{
	Name: "-s", Desc: "Serves rpc handler API via apphost interface.",
	Func: ServeFunc,
}

func ServeFunc(ctx context.Context, root *cmd.Root) error {
	handler := cmd.Handler(*root)
	r := Default().Router(handler)
	return r.Run(ctx)
}
