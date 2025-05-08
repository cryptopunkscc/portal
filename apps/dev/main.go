package main

import (
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/target/dev/broadcast"
)

func main() { cli.Run(Application{}.cliHandler()) }

type Application struct{}

func (a Application) cliHandler() cmd.Handler {
	return cmd.Handler{
		Func: apphost.Default.Rpc().Router(a.netHandler()).Run,
		Name: "dev",
		Desc: "Portal development service.",
	}
}

func (a Application) netHandler() cmd.Handler {
	return cmd.Handler{
		Name: "dev.portal",
		Sub: cmd.Handlers{{
			Func: broadcast.New().BroadcastMsg,
			Name: "broadcast",
		}},
	}
}
