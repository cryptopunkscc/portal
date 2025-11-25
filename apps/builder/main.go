package main

import (
	"context"

	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/target/all"
	"github.com/cryptopunkscc/portal/target/project"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: build,
		Name: "portal-build",
		Desc: "Builds portal project and generates application bundle.",
		Params: cmd.Params{
			{Name: "clean c", Type: "bool", Desc: "Clean build directory before building."},
			{Name: "pack p", Type: "bool", Desc: "Pack app bundle after successful build."},
			{Name: "out o", Type: "string", Desc: "optional Path to output directory."},
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

type buildOpts struct {
	Clean bool   `cli:"clean c"`
	Pack  bool   `cli:"pack p"`
	Dir   string `cli:"out o"`
}

func build(ctx context.Context, ops buildOpts, src string) error {
	var args []string
	if ops.Pack {
		args = append(args, "pack")
	}
	if ops.Clean {
		args = append(args, "clean")
	}
	if ops.Dir != "" {
		args = append(args, ops.Dir)
	}
	return all.BuildRecursive(ctx, src, args...)
}
