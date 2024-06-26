package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	manifest "github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/build"
	"github.com/cryptopunkscc/portal/feat/create"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/version"
	osExec "github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	create2 "github.com/cryptopunkscc/portal/runner/create"
	"github.com/cryptopunkscc/portal/runner/dev"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/go_dev"
	"github.com/cryptopunkscc/portal/runner/pack"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/msg"
	"github.com/cryptopunkscc/portal/target/portals"
	"os"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	println("...")
	plog.ErrorStackTrace = true
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")

	portalPort.InitPrefix("dev")
	port := target.Port{Base: "portal"}
	portOpen := port.Route("open")
	portMsg := port.Route("broadcast")
	scope := feature.Scope[target.Portal]{
		Astral:         serve.CheckAstral,
		Executable:     "portal-dev",
		Port:           port,
		WrapApi:        NewAdapter,
		WaitGroup:      &sync.WaitGroup{},
		TargetCache:    target.NewCache[target.Portal](),
		NewRunService:  service.NewRun,
		TargetFinder:   portals.NewFind,
		NewExecTarget:  exec.NewRun[target.Portal],
		GetPath:        featApps.Path,
		FeatObserve:    featApps.Observe,
		JoinTarget:     query.NewRunner[target.App](portOpen).Run,
		DispatchTarget: query.NewRunner[target.App](portOpen).Start,
		Processes:      &sig.Map[string, target.Portal]{},
	}
	scope.RpcHandlers = rpc.Handlers{
		portMsg.Name: msg.NewBroadcast(portMsg, scope.GetProcesses()).BroadcastMsg,
	}
	scope.DispatchService = scope.GetServeFeature().Dispatch

	featBuild := build.NewFeat(dist.NewRun, pack.Run)
	featCreate := create.NewFeat(create2.NewRun, featBuild.Dist).Run

	goRunner := go_dev.NewRunner(
		featBuild.Dist,
		portMsg,
		func(ctx context.Context, src target.DistExec) (err error) { return scope.GetExecTarget()(ctx, src) },
	)
	scope.NewRunTarget = dev.NewRun(portMsg, goRunner)

	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(scope.GetDispatchFeature())
	cli.Open(scope.GetOpenFeature())
	cli.Create(create.List, featCreate)
	cli.Build(featBuild.Run)
	cli.Portals(scope.GetTargetFind())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	scope.WaitGroup.Wait()
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }
