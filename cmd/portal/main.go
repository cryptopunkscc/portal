package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/clir"
	feature "github.com/cryptopunkscc/go-astral-js/feat"
	featApps "github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	osExec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/runner/app"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/runner/service"
	"github.com/cryptopunkscc/go-astral-js/runner/tray"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"os"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	println("...")
	log := plog.New().I().Set(&ctx).Scope("main")
	log.Println("starting portal", os.Args)
	defer log.Println("closing portal")

	scope := feature.Scope[target.App]{
		Port:            "portal",
		WrapApi:         NewAdapter,
		WaitGroup:       &sync.WaitGroup{},
		TargetCache:     target.NewCache[target.App](),
		NewRunTarget:    app.NewRun,
		NewRunTray:      tray.NewRun,
		NewRunService:   service.NewRun,
		ExecTarget:      exec.NewRun[target.App]("portal"),
		TargetFinder:    apps.NewFind,
		GetPath:         featApps.Path,
		FeatObserve:     featApps.Observe,
		JoinTarget:      query.NewRunner[target.App]("portal.open").Run,
		DispatchTarget:  query.NewRunner[target.App]("portal.open").Start,
		DispatchService: exec.NewService("portal").Start,
	}
	scope.RpcHandlers = rpc.Handlers{
		"install":   featApps.Install,
		"uninstall": featApps.Uninstall,
	}

	cli := clir.NewCli(ctx, manifest.Name, manifest.Description, version.Run)
	cli.Dispatch(scope.GetDispatchFeature())
	cli.Serve(scope.GetServeFeature().Run)
	cli.Open(scope.GetOpenFeature())
	cli.Apps(scope.GetTargetFind())
	cli.Install(featApps.Install)
	cli.Uninstall(featApps.Uninstall)
	cli.List(featApps.List)

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	scope.WaitGroup.Wait()
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }
