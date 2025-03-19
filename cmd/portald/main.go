package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/rpc/cli"
	"github.com/cryptopunkscc/portal/core/rpc/cmd"
	"github.com/cryptopunkscc/portal/pkg/plog"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
)

var application = &Application[Portal_]{}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	application.Shutdown = cancel
	log := plog.New().D().Scope("portald").Set(&ctx)
	go singal.OnShutdown(log, cancel)
	handler := application.handler()
	cmd.InjectHelp(&handler)
	err := cli.New(handler).Run(ctx)
	if err != nil {
		log.E().Println("finished with error:", err)
	}
}
