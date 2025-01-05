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
	go sig.OnShutdown(cancel)
	log := plog.New().D().Scope(handler.Name).Set(&ctx)
	err := cli.New(handler).Run(ctx)
	if err != nil {
		log.Println(err)
	}
	cancel()
}
