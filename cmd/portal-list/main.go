package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	find "github.com/cryptopunkscc/portal/factory/find/portal"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)
	cli := clir.NewCli(ctx,
		"portal-build",
		"Builds portal project and generates application bundle.",
		version.Run)
	cli.Portals(find.Create[target.Portal_]())
	if err := cli.Run(); err != nil {
		panic(err)
	}
}
