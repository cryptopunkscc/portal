package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/factory/dev"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/goja_pro"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{ dev.Module[PortalJs] }

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: open.Runner[PortalJs](&a),
		Name: "dev-js",
		Desc: "Start portal js app development in goja runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}
}

func (a Application) Runner() Run[PortalJs] {
	return multi.Runner[PortalJs](
		reload.Mutable(bind.BackendRuntime(), PortMsg, goja_pro.NewRunner),
		reload.Mutable(bind.BackendRuntime(), PortMsg, goja_dist.NewRunner),
		reload.Immutable(bind.BackendRuntime(), PortMsg, goja.NewRunner),
	)
}
