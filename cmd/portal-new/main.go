package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/create"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/template"
	"log"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)
	cli := clir.NewCli(ctx,
		"portal-new",
		"Create new portal project from template.",
		version.Run)
	cli.Create(
		template.List,
		create.Create())
	if err := cli.Run(); err != nil {
		panic(err)
	}
	log.Println("* done")
}
