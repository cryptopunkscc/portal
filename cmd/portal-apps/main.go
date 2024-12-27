package main

import (
	"context"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)

	err := cli.New(cmd.Handler{
		Name: "portal-apps",
		Desc: "Portal applications management.",
		Sub: cmd.Handlers{
			{
				Func: appstore.ListApps,
				Name: "list l",
				Desc: "List installed apps.",
			},
			{
				Func: appstore.Install,
				Name: "install i",
				Desc: "Install app from a given portal app bundle path.",
				Params: cmd.Params{
					{Type: "string", Desc: "Path to containing directory"},
				},
			},
			{
				Func: appstore.Uninstall,
				Name: "delete d",
				Desc: "Uninstall app.",
				Params: cmd.Params{
					{Type: "string", Desc: "Application name or package name"},
				},
			},
		},
	}).Run(ctx)

	if err != nil {
		panic(err)
	}
}
