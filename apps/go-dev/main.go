package main

import (
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/go_project"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/target/source"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: source.File.NewRun(go_project.Runner().Try),
		Name: "dev-go",
		Desc: "Start portal golang app development.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}
