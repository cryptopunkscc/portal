package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)
	cli := clir.NewCli(ctx,
		"Portal-apps",
		"Portal applications management.",
		version.Run)
	cli.List(appstore.ListApps)
	cli.Install(appstore.Install)
	cli.Uninstall(appstore.Uninstall)
	if err := cli.Run(); err != nil {
		panic(err)
	}
}
