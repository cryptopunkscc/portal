package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails_dev"
	"github.com/cryptopunkscc/portal/runner/wails_dist"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/portals"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-wails").Set(&ctx)

	go sig.OnShutdown(cancel)
	portalPort.InitPrefix("dev")

	scope := feature.Scope[target.PortalHtml]{
		WrapApi:        NewAdapter,
		GetPath:        featApps.Path,
		TargetFinder:   portals.NewFind[target.PortalHtml],
		TargetCache:    target.NewCache[target.PortalHtml](),
		DispatchTarget: query.NewRunner[target.PortalHtml](target.PortOpen).Start,
	}
	scope.NewRunTarget = func(newApi target.NewApi) target.Run[target.PortalHtml] {
		return multi.NewRunner[target.PortalHtml](
			reload.Immutable(newApi, target.PortMsg, wails_dev.NewRunner), // FIXME propagate sendMsg
			reload.Mutable(newApi, target.PortMsg, wails_dist.NewRunner),
			reload.Immutable(newApi, target.PortMsg, wails.NewRunner),
		).Run
	}

	cli := clir.NewCli(ctx,
		"Portal-dev-wails",
		"Portal html development driven by wails.",
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
