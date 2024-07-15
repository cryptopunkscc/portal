package main

import (
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	osExec "github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dev"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/portals"
	"golang.org/x/net/context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	log := plog.New().D().Scope("dev-goja").Set(&ctx)
	portalPort.InitPrefix("dev")
	port := target.NewPort("portal")
	portOpen := port.Route("open")
	portMsg := port.Route("broadcast")

	scope := feature.Scope[target.PortalJs]{
		WrapApi:        NewAdapter,
		GetPath:        featApps.Path,
		TargetFinder:   portals.NewFind[target.PortalJs],
		TargetCache:    target.NewCache[target.PortalJs](),
		DispatchTarget: query.NewRunner[target.PortalJs](portOpen).Start,
	}
	scope.NewRunTarget = func(newApi target.NewApi) target.Run[target.PortalJs] {
		return multi.NewRunner[target.PortalJs](
			reload.Mutable(newApi, portMsg, goja_dev.NewRunner),
			reload.Mutable(newApi, portMsg, goja_dist.NewRunner),
			reload.Immutable(newApi, portMsg, goja.NewRunner),
		).Run
	}

	cli := clir.NewCli(ctx,
		"Portal-dev-goja",
		"Portal js development driven by goja.",
		version.Run,
	)
	cli.Open(scope.GetOpenFeature())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }
