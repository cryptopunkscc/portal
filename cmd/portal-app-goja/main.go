package main

import (
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	osExec "github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"golang.org/x/net/context"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

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
