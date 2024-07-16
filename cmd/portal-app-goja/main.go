package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-goja").Set(&ctx)

	go sig.OnShutdown(cancel)

	scope := feature.Scope[target.AppJs]{
		WrapApi:        NewAdapter,
		GetPath:        featApps.Path,
		TargetFinder:   apps.NewFind[target.AppJs],
		TargetCache:    target.NewCache[target.AppJs](),
		DispatchTarget: query.NewRunner[target.AppJs](target.PortOpen).Start,
		NewRunTarget:   goja.NewRun,
	}

	cli := clir.NewCli(ctx,
		"Portal-goja",
		"Portal js runner driven by goja.",
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
