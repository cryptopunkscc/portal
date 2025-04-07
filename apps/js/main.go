package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/version"
)

func main() { cli.Run(Application{}.Handler()) }

type Application struct{}

func (a Application) Handler() cmd.Handler {
	return cmd.Handler{
		Func: open.NewRun[AppJs](&a),
		Name: "js",
		Desc: "Start portal app in goja runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}
func (a Application) Runner() Run[AppJs]       { return goja.NewRun(bind.NewBackendCore) }
func (a Application) Resolver() Resolve[AppJs] { return apps.Resolver[AppJs]() }
