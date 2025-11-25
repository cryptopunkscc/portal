package cli

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd/help"
	"github.com/cryptopunkscc/portal/pkg/sig"
)

func Run(handler cmd.Handler) {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope(handler.Name).Set(&ctx)
	go sig.OnShutdown(log, cancel)
	if !help.Check(handler) {
		help.Inject(&handler)
	}
	err := New(handler).Run(ctx)
	if err != nil {
		log.Println(err)
	}
	cancel()
}
