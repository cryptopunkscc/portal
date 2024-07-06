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
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/go_dev"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dev"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/pack"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/runner/template"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails_dev"
	"github.com/cryptopunkscc/portal/runner/wails_dist"
	"github.com/cryptopunkscc/portal/target"
	js "github.com/cryptopunkscc/portal/target/js/embed"
	"github.com/cryptopunkscc/portal/target/msg"
	"github.com/cryptopunkscc/portal/target/portals"
	"github.com/cryptopunkscc/portal/target/sources"
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
	scope.NewRunTarget = func(newApi target.NewApi) target.Run[target.Portal] {
		return multi.NewRunner(
			reload.Mutable(newApi, portMsg, goja_dev.NewRunner),
			reload.Immutable(newApi, portMsg, wails_dev.NewRunner), // FIXME propagate sendMsg
			reload.Mutable(newApi, portMsg, go_dev.NewAdapter(scope.GetExecTarget())),
			reload.Mutable(newApi, portMsg, goja_dist.NewRunner),
			reload.Mutable(newApi, portMsg, wails_dist.NewRunner),
			reload.Immutable(newApi, portMsg, goja.NewRunner),
			reload.Immutable(newApi, portMsg, wails.NewRunner),
		).Run
	}
	scope.DispatchService = scope.GetServeFeature().Dispatch

	featBuild := build.NewFeat(
		dist.NewRun, pack.Run,
		sources.FromFS[target.NodeModule](js.PortalLibFS),
	)
	featCreate := create.NewFeat(template.NewRun, featBuild.Dist)

	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(scope.GetDispatchFeature())
	cli.Open(scope.GetOpenFeature())
	cli.Create(template.List, featCreate.Run)
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
