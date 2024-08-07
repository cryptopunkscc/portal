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
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	signal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/all_build"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/pack"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/template"
	. "github.com/cryptopunkscc/portal/target"
	js "github.com/cryptopunkscc/portal/target/js/embed"
	"github.com/cryptopunkscc/portal/target/msg"
	"os"
)

func main() {
	mod := Module[Portal_]{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	mod.CancelFunc = cancel
	go signal.OnShutdown(cancel)
	println("...")
	plog.ErrorStackTrace = true
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(mod.FeatDev())
	cli.Create(template.List, mod.FeatCreate().Run)
	cli.Build(mod.FeatBuild().Run)
	cli.Portals(mod.TargetFind())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	mod.WaitGroup().Wait()
}

type Module[T Portal_] struct{ srv.Module[T] }

func (d *Module[T]) Executable() string   { return "portal-dev" }
func (d *Module[T]) Astral() serve.Astral { return serve.CheckAstral }
func (d *Module[T]) RpcHandlers() rpc.Handlers {
	return rpc.Handlers{
		PortMsg.Name: msg.NewBroadcast(PortMsg, d.Processes()).BroadcastMsg,
	}
}
func (d *Module[T]) TargetResolve() Resolve[T] { return sources.Resolver[T]() }
func (d *Module[T]) TargetRun() Run[T] {
	return multi.NewRunner[T](
		app.Run(exec.NewPortal[PortalJs]("portal-dev-goja", "o").Run),
		app.Run(exec.NewPortal[PortalHtml]("portal-dev-wails", "o").Run),
		app.Run(exec.NewPortal[ProjectGo]("portal-dev-go", "o").Run),
		app.Run(exec.NewPortal[AppExec]("portal-dev-exec", "o").Run),
	).Run
}
func (d *Module[T]) Priority() Priority {
	return []Matcher{
		Match[Project_],
		Match[Dist_],
		Match[Bundle_],
	}
}
func (d *Module[T]) JoinTarget() Dispatch      { return query.NewOpen().Run }
func (d *Module[T]) DispatchService() Dispatch { return serve.Inject[T](d).Dispatch }
func (d *Module[T]) FeatDev() Dispatch         { return dispatch.Inject(d).Run }
func (d *Module[T]) FeatCreate() *create.Feat {
	return create.NewFeat(template.NewRun, d.FeatBuild().Dist)
}
func (d *Module[T]) FeatBuild() *build.Feat[T] {
	return build.NewFeat(
		d.TargetResolve(),
		all_build.NewRun,
		pack.Run,
		List(
			Any[NodeModule](
				Skip("node_modules"),
				Try(npm.Resolve),
			),
			source.Embed(js.PortalLibFS),
		),
	)
}
