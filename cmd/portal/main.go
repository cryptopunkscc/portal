package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/builder"
	"github.com/cryptopunkscc/go-astral-js/clir"
	featApps "github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	osexec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/runner/app"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"github.com/cryptopunkscc/go-astral-js/runner/tray"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"os"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osexec.OnShutdown(cancel)

	println("...")
	log := plog.New().I().Set(&ctx)
	log.Scope("main").Println("starting portal", os.Args)

	scope := builder.Scope[target.App]{
		Port:            "portal",
		WrapApi:         NewAdapter,
		WaitGroup:       &sync.WaitGroup{},
		TargetCache:     target.NewCache[target.App](),
		NewTargetRun:    app.NewRun,
		NewTray:         tray.New,
		NewServe:        serve.NewRun,
		TargetFinder:    apps.NewFind,
		ExecTarget:      exec.NewRun[target.App]("portal"),
		AppsPath:        featApps.Path,
		FeatObserve:     featApps.Observe,
		FeatInstall:     featApps.Install,
		FeatUninstall:   featApps.Uninstall,
		DispatchTarget:  query.NewRunner[target.App]("portal.open").Run,
		DispatchService: exec.NewService("portal").Run,
	}

	cli := clir.NewCli(ctx, manifest.Name, manifest.Description, version.Run)
	cli.Dispatch(scope.GetDispatchFeature())
	cli.Serve(scope.GetServeFeature())
	cli.Open(scope.GetOpenFeature())
	cli.Install(scope.FeatInstall)
	cli.Uninstall(scope.FeatUninstall)
	cli.Apps(scope.GetTargetFind())
	cli.List(featApps.List)

	err := cli.Run()
	cancel()
	if err != nil {
		log.Println(err)
	}
	scope.WaitGroup.Wait()
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }
