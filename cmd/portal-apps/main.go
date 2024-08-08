package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)
	cli := clir.NewCli(ctx,
		"Portal-apps",
		"Portal applications management.",
		version.Run)
	cli.List(apps.List)
	cli.Install(apps.Install)
	cli.Uninstall(apps.Uninstall)
	if err := cli.Run(); err != nil {
		panic(err)
	}
}
