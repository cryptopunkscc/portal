package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	create "github.com/cryptopunkscc/portal/factory/apps"
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

	apps := create.Default()
	cli.List(apps.List)
	cli.Install(apps.InstallFromPath)
	cli.Uninstall(apps.Uninstall)
	if err := cli.Run(); err != nil {
		panic(err)
	}
}
