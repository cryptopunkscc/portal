package main

import (
	"github.com/cryptopunkscc/portal/factory/build"
	"github.com/cryptopunkscc/portal/runner/clean"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: build.Create().Run,
		Name: "portal-build",
		Desc: "Builds portal project and generates application bundle.",
		Params: cmd.Params{
			{Type: "string", Desc: "Path to project directory. Takes current directory as default."},
		},
		Sub: cmd.Handlers{
			{
				Func: clean.Runner(),
				Name: "clean c",
				Desc: "Clean target directories from build artifacts without building.",
				Params: cmd.Params{
					{Type: "string", Desc: "Path to project directory. Default is '.'"},
				},
			},
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}
