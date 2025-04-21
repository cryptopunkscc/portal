package main

import (
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/version"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	r := exec.DefaultRunner()
	return cmd.Handler{
		Func: source.File.NewRun(
			r.Bundle().Try,
			r.Dist().Try,
		),
		Name: "dev-exec",
		Desc: "Portal development runner for executables.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}
