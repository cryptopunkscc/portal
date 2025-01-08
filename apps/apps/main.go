package main

import (
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/install"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() { cli.Run(Application{}.Handler()) }

type Application struct{}

func (a Application) Handler() cmd.Handler {
	return cmd.Handler{
		Name: "apps",
		Desc: "Portal applications management.",
		Sub: cmd.Handlers{
			{
				Func: appstore.ListApps,
				Name: "list l",
				Desc: "List installed apps.",
			},
			{
				Func: install.Runner(a.dir()).Run,
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
	}
}

func (a Application) dir() string { return PortalAppsDir }
