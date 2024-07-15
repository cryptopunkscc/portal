package main

import (
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	osExec "github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/runner/exec"
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
		Port:            target.PortPortal,
		GetPath:         featApps.Path,
		TargetFinder:    apps.NewFind[target.AppJs],
		TargetCache:     target.NewCache[target.AppJs](),
		DispatchService: exec.NewDispatcher("portal-app").Dispatch,
		JoinTarget:      query.NewRunner[target.App](target.PortOpen).Run,
		NewRunTarget:    goja.NewRun,
	}

	cli := clir.NewCli(ctx,
		"Portal",
		"Portal command line.",
		version.Run,
	)
	cli.Dispatch(scope.GetDispatchFeature())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}
