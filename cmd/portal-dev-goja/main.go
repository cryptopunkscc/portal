package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/di/run/dev"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dev"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/reload"
	. "github.com/cryptopunkscc/portal/target"
)

func main() {
	mod := Module{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-goja").Set(&ctx)
	go sig.OnShutdown(cancel)
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx,
		"Portal-dev-goja",
		"Portal js development driven by goja.",
		version.Run,
	)
	cli.Open(mod.FeatOpen())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ dev.Module[PortalJs] }

func (d *Module) NewRunTarget(newApi NewApi) Run[PortalJs] {
	return multi.NewRunner[PortalJs](
		reload.Mutable(newApi, PortMsg, goja_dev.NewRunner),
		reload.Mutable(newApi, PortMsg, goja_dist.NewRunner),
		reload.Immutable(newApi, PortMsg, goja.NewRunner),
	).Run
}
