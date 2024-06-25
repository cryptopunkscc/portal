package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	manifest "github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	osExec "github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/runner/tray"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"os"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	println("...")

	plog.ErrorStackTrace = true
	log := plog.New().I().Set(&ctx).Scope("main")
	log.Println("starting portal", os.Args)
	defer log.Println("closing portal")

	port := target.Port{Base: "portal"}
	portOpen := port.Route("open")
	executable := "portal"
	scope := feature.Scope[target.App]{
		Astral:          exec.Astral,
		Executable:      executable,
		Port:            port,
		WrapApi:         NewAdapter,
		WaitGroup:       &sync.WaitGroup{},
		TargetCache:     target.NewCache[target.App](),
		NewRunTarget:    app.NewRun,
		NewRunTray:      tray.NewRun,
		NewRunService:   service.NewRun,
		NewExecTarget:   exec.NewRun[target.App],
		TargetFinder:    apps.NewFind,
		GetPath:         featApps.Path,
		FeatObserve:     featApps.Observe,
		JoinTarget:      query.NewRunner[target.App](portOpen).Run,
		DispatchTarget:  query.NewRunner[target.App](portOpen).Start,
		DispatchService: exec.NewDispatcher(executable).Dispatch,
		Processes:       &sig.Map[string, target.App]{},
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
