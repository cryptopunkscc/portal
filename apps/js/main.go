package main

import (
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/target/source"
)

func main() { cli.Run(Application{}.Handler()) }

type Application struct{}

func (a Application) Handler() cmd.Handler {
	return cmd.Handler{
		Func: source.File.NewRun(goja.Runner(bind.NewBackendCore).Try),
		Name: "js",
		Desc: "Start portal app in goja runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Name},
		},
	}
}
