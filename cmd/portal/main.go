package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/feat/dispatch"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/target"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app").Set(&ctx)

	go sig.OnShutdown(cancel)

	dispatchFeat := dispatch.NewFeat(
		target.PortPortal,
		query.NewRunner[target.App](target.PortOpen).Run,
		exec.NewDispatcher("portal-app").Dispatch,
	)

	cli := clir.NewCli(ctx,
		"Portal",
		"Portal command line.",
		version.Run,
	)
	cli.Dispatch(dispatchFeat)

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}
