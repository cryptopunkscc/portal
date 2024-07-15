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
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/pack"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/runner/template"
	"github.com/cryptopunkscc/portal/target"
	js "github.com/cryptopunkscc/portal/target/js/embed"
	"github.com/cryptopunkscc/portal/target/msg"
	"github.com/cryptopunkscc/portal/target/portals"
	"github.com/cryptopunkscc/portal/target/sources"
	"os"
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

	scope := feature.Scope[target.Portal]{
		Astral:        serve.CheckAstral,
		Executable:    "portal-dev",
		Port:          target.PortPortal,
		TargetCache:   target.NewCache[target.Portal](),
		NewRunService: service.NewRun,
		TargetFinder:  portals.NewFind[target.Portal],
		GetPath:       featApps.Path,
		FeatObserve:   featApps.Observe,
		JoinTarget:    query.NewRunner[target.App](target.PortOpen).Run,
		Processes:     &sig.Map[string, target.Portal]{},
		NewExecTarget: func(_ string, _ string) target.Run[target.Portal] {
			return multi.NewRunner[target.Portal](
				app.Run(exec.NewPortal[target.PortalJs]("portal-dev-goja", "o").Run),
				app.Run(exec.NewPortal[target.PortalHtml]("portal-dev-wails", "o").Run),
				app.Run(exec.NewPortal[target.ProjectGo]("portal-dev-go", "o").Run),
				app.Run(exec.NewPortal[target.AppExec]("portal-dev-exec", "o").Run),
			).Run
		},
	}
	scope.RpcHandlers = rpc.Handlers{
		target.PortMsg.Name: msg.NewBroadcast(target.PortMsg, scope.GetProcesses()).BroadcastMsg,
	}
	scope.DispatchService = scope.GetServeFeature().Dispatch

	featBuild := build.NewFeat(
		dist.NewRun, pack.Run,
		sources.FromFS[target.NodeModule](js.PortalLibFS),
	)
	featCreate := create.NewFeat(template.NewRun, featBuild.Dist)

	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(scope.GetDispatchFeature())
	cli.Create(template.List, featCreate.Run)
	cli.Build(featBuild.Run)
	cli.Portals(scope.GetTargetFind())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	scope.WaitGroup.Wait()
}
