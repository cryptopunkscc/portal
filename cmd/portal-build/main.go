package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/build"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/clean"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)
	cli := clir.NewCli(ctx,
		"portal-build",
		"Builds portal project and generates application bundle.",
		version.Run)
	cli.Build(build.Create().Run, clean.NewRunner().Call)
	if err := cli.Run(); err != nil {
		panic(err)
	}
}
