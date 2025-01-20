package main

import (
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/dev"
	"github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() { cli.Run(Application{}.cliHandler()) }

type Application struct{}

func (a Application) cliHandler() cmd.Handler {
	return cmd.Handler{
		Func: apphost.NewRouter(a.netHandler()).Run,
		Name: "dev",
		Desc: "Portal development service.",
	}
}

func (a Application) netHandler() cmd.Handler {
	return cmd.Handler{
		Name: "dev.portal",
		Sub: cmd.Handlers{{
			Func: dev.NewBroadcast().BroadcastMsg,
			Name: "broadcast",
		}},
	}
}
