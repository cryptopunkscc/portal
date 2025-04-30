package main

import (
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/goja_pro"
	"github.com/cryptopunkscc/portal/target/source"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: source.File.NewRun(
			goja_pro.Runner(bind.NewBackendCore).Try,
			goja_dist.Runner(bind.NewBackendCore).Try,
			goja.Runner(bind.NewBackendCore).Try,
		),
		Name: "dev-js",
		Desc: "Start portal js app development in goja runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Name},
		},
	}
}
