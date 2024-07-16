package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-wails").Set(&ctx)

	go sig.OnShutdown(cancel)

	scope := feature.Scope[target.AppHtml]{
		WrapApi:        NewAdapter,
		GetPath:        featApps.Path,
		TargetFinder:   apps.NewFind[target.AppHtml],
		TargetCache:    target.NewCache[target.AppHtml](),
		DispatchTarget: query.NewRunner[target.AppHtml](target.PortOpen).Start,
		NewRunTarget:   wails.NewRun,
	}

	cli := clir.NewCli(ctx,
		"Portal-wails",
		"Portal html runner driven by wails.",
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
