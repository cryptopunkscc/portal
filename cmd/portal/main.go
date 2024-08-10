package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/feat/start"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/request/exec"
	"github.com/cryptopunkscc/portal/request/query"
	"github.com/cryptopunkscc/portal/target"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app").Set(&ctx)

	go sig.OnShutdown(cancel)

	cli := clir.NewCli(ctx,
		"Portal",
		"Portal command line.",
		version.Run,
	)
	cli.Request(start.Feat(deps{}))

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type deps struct{}

func (m deps) Port() target.Port          { return target.PortPortal }
func (m deps) Serve() target.Request      { return exec.Request("portal-app") }
func (m deps) JoinTarget() target.Request { return query.Request.Run }
