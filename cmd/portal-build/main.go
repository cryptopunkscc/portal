package main

import (
	"context"
	"github.com/cryptopunkscc/portal/factory/build"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/clean"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"log"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)

	err := cli.New(cmd.Handler{
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
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}).Run(ctx)

	if err != nil {
		panic(err)
	}
	log.Println("* done")
}
