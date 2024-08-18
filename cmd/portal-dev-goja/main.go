package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/factory/dev"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/goja_pro"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/reload"
)

func main() {
	mod := &Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-goja").Set(&ctx)
	go sig.OnShutdown(cancel)
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx,
		"Portal-dev-goja",
		"Portal js development driven by goja.",
		version.Run,
	)
	cli.Open(open.Feat[PortalJs](mod))
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ dev.Module[PortalJs] }

func (d *Module) Runner() Run[PortalJs] {
	return multi.Runner[PortalJs](
		reload.Mutable(bind.BackendRuntime(), PortMsg, goja_pro.NewRunner),
		reload.Mutable(bind.BackendRuntime(), PortMsg, goja_dist.NewRunner),
		reload.Immutable(bind.BackendRuntime(), PortMsg, goja.NewRunner),
	)
}
