package main

import (
	"context"
	manifest "github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/di/srv"
	"github.com/cryptopunkscc/portal/feat/build"
	"github.com/cryptopunkscc/portal/feat/create"
	"github.com/cryptopunkscc/portal/feat/dispatch"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	signal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/pack"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/template"
	. "github.com/cryptopunkscc/portal/target"
	js "github.com/cryptopunkscc/portal/target/js/embed"
	"github.com/cryptopunkscc/portal/target/msg"
	"github.com/cryptopunkscc/portal/target/portals"
	"github.com/cryptopunkscc/portal/target/sources"
	"os"
)

func main() {
	mod := Module[Portal]{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	go signal.OnShutdown(cancel)
	println("...")
	plog.ErrorStackTrace = true
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(mod.FeatDispatch())
	cli.Create(template.List, mod.FeatCreate().Run)
	cli.Build(mod.FeatBuild().Run)
	cli.Portals(mod.TargetFind())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	mod.WaitGroup().Wait()
}

type Module[T Portal] struct{ srv.Module[T] }

func (d *Module[T]) Executable() string        { return "portal-dev" }
func (d *Module[T]) Astral() serve.Astral      { return serve.CheckAstral }
func (d *Module[T]) JoinTarget() Dispatch      { return query.NewRunner[Portal](PortOpen).Run }
func (d *Module[T]) TargetFinder() Finder[T]   { return portals.NewFind[T] }
func (d *Module[T]) DispatchService() Dispatch { return di.Single(serve.Inject[T], d).Dispatch }
func (d *Module[T]) FeatDispatch() Dispatch    { return di.Single(dispatch.Inject, d) }
func (d *Module[T]) FeatCreate() *create.Feat {
	return create.NewFeat(template.NewRun, d.FeatBuild().Dist)
}
func (d *Module[T]) FeatBuild() *build.Feat {
	return build.NewFeat(
		dist.NewRun, pack.Run,
		sources.FromFS[NodeModule](js.PortalLibFS),
	)
}
func (d *Module[T]) RpcHandlers() rpc.Handlers {
	return rpc.Handlers{
		PortMsg.Name: msg.NewBroadcast(PortMsg, d.Processes()).BroadcastMsg,
	}
}
func (d *Module[T]) TargetRun() Run[T] {
	return multi.NewRunner[T](
		app.Run(exec.NewPortal[PortalJs]("portal-dev-goja", "o").Run),
		app.Run(exec.NewPortal[PortalHtml]("portal-dev-wails", "o").Run),
		app.Run(exec.NewPortal[ProjectGo]("portal-dev-go", "o").Run),
		app.Run(exec.NewPortal[AppExec]("portal-dev-exec", "o").Run),
	).Run
}
