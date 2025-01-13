package main

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/install"
	"github.com/cryptopunkscc/portal/runner/observe"
	"github.com/cryptopunkscc/portal/runner/uninstall"
	apps2 "github.com/cryptopunkscc/portal/runtime/apps"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() { cli.Run(Application{}.cliHandler()) }

type Application struct{}

func (a Application) cliHandler() cmd.Handler {
	return cmd.Handler{
		Name: "apps",
		Desc: "Portal applications management.",
		Sub: cmd.Handlers{
			{
				Func: apps.ResolveAll.List(a.src()),
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
				Func: uninstall.Runner(a.src()),
				Name: "delete d",
				Desc: "Uninstall app.",
				Params: cmd.Params{
					{Type: "string", Desc: "Application name or package name"},
				},
			},
			a.apiHandler(),
		},
	}
}

func (a Application) apiHandler() cmd.Handler {
	return cmd.Handler{
		Name: "serve s",
		Desc: "Serve apps management.",
		Sub: cmd.Handlers{
			{
				Name: "observe",
				Func: observe.NewRun(a.dir()),
			},
		},
	}
}

func (a Application) dir() string        { return apps2.Dir }
func (a Application) src() target.Source { return apps2.Source }
