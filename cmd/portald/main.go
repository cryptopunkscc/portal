package main

import (
	"context"
	_ "github.com/cryptopunkscc/portal/api/env/desktop"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd/help"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
)

var application = &Application{}

func init() {
	plog.Default = plog.New().Scope("portald")
	application.ExtraTokens = []string{"portal"}
}

func main() {
	ctx := context.Background()
	log := plog.Default.D().Set(&ctx)
	go singal.OnShutdown(log, application.Stop)

	c := application.commands()
	help.Inject(&c)
	err := cli.New(c).Run(ctx)
	if err != nil {
		log.E().Println("finished with error:", err)
	}
}
