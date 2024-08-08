package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/dispatch/exec"
	"github.com/cryptopunkscc/portal/dispatch/query"
	"github.com/cryptopunkscc/portal/feat/start"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/target"
)

func main() {
	mod := Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app").Set(&ctx)

	go sig.OnShutdown(cancel)

	cli := clir.NewCli(ctx,
		"Portal",
		"Portal command line.",
		version.Run,
	)
	cli.Dispatch(mod.FeatDispatch())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ joinTarget target.Dispatch }

func (m Module) Port() target.Port                { return target.PortPortal }
func (m Module) DispatchService() target.Dispatch { return exec.NewDispatcher("portal-app").Dispatch }
func (m Module) JoinTarget() target.Dispatch      { return m.joinTarget }
func (m Module) FeatDispatch() target.Dispatch {
	m.joinTarget = query.NewOpen().Run
	return start.Inject(m).Run
}
