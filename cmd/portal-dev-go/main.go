package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/go_dev"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/portals"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-go").Set(&ctx)

	go sig.OnShutdown(cancel)
	port.InitPrefix("dev")

	scope := feature.Scope[target.ProjectGo]{
		WrapApi:        NewAdapter,
		GetPath:        featApps.Path,
		TargetFinder:   portals.NewFind[target.ProjectGo],
		TargetCache:    target.NewCache[target.ProjectGo](),
		DispatchTarget: query.NewRunner[target.ProjectGo](target.PortOpen).Start,
	}
	scope.NewRunTarget = func(newApi target.NewApi) target.Run[target.ProjectGo] {
		return multi.NewRunner[target.ProjectGo](
			reload.Mutable(newApi, target.PortMsg, go_dev.NewAdapter(app.Run(exec.NewDistRunner().Run))),
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
