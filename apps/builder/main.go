package main

import (
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/target/all"
	"github.com/cryptopunkscc/portal/target/project"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: all.BuildRecursive,
		Name: "portal-build",
		Desc: "Builds portal project and generates application bundle.",
		Params: cmd.Params{
			{Type: "string", Desc: "Path to project directory. Takes current directory as default."},
		},
		Sub: cmd.Handlers{
			{
				Func: project.Cleaner(),
				Name: "clean c",
				Desc: "Clean target directories from build artifacts without building.",
				Params: cmd.Params{
					{Type: "string", Desc: "Path to project directory. Default is '.'"},
				},
			},
			{Name: "v", Desc: "Print version.", Func: version.Name},
		},
	}
}
