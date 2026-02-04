package main

import (
	"github.com/cryptopunkscc/portal/pkg/util/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/util/rpc/cmd"
)

func main() { cli.Run(handler) }

var handler = cmd.Handler{
	Name: "portal-tools",
	Sub: cmd.Handlers{
		ListTargetsHandler,
		ListGoImportsHandler,
	},
}
