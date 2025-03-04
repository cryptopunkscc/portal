package cli

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func Run(handler cmd.Handler) {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope(handler.Name).Set(&ctx)
	go sig.OnShutdown(log, cancel)
	if !cmd.HasHelp(handler) {
		cmd.InjectHelp(&handler)
	}
	err := cli.New(handler).Run(ctx)
	if err != nil {
		log.Println(err)
	}
	cancel()
}
