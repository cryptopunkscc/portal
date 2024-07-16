package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	manifest "github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/clir"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
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
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/runner/template"
	"github.com/cryptopunkscc/portal/runner/tray"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/find"
	js "github.com/cryptopunkscc/portal/target/js/embed"
	"github.com/cryptopunkscc/portal/target/msg"
	"github.com/cryptopunkscc/portal/target/portals"
	"github.com/cryptopunkscc/portal/target/sources"
	"github.com/cryptopunkscc/portal/target/spawn"
	"os"
	"sync"
)

func main() {
	deps := NewDeps[Portal]()
	ctx, cancel := context.WithCancel(context.Background())
	go signal.OnShutdown(cancel)
	println("...")
	plog.ErrorStackTrace = true
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(deps.FeatDispatch())
	cli.Create(template.List, deps.FeatCreate().Run)
	cli.Build(deps.FeatBuild().Run)
	cli.Portals(deps.TargetFind())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	deps.WaitGroup().Wait()
}

type Deps[T Portal] struct {
	di.Cache
	wg        sync.WaitGroup
	processes sig.Map[string, T]
	cache     *Cache[T]
}

func NewDeps[T Portal]() *Deps[T]                 { return &Deps[T]{cache: NewCache[T]()} }
func (d *Deps[T]) Path() Path                     { return featApps.Path }
func (d *Deps[T]) Port() Port                     { return PortPortal }
func (d *Deps[T]) Executable() string             { return "portal-dev" }
func (d *Deps[T]) WaitGroup() *sync.WaitGroup     { return &d.wg }
func (d *Deps[T]) Processes() *sig.Map[string, T] { return &d.processes }
func (d *Deps[T]) TargetCache() *Cache[T]         { return d.cache }
func (d *Deps[T]) TargetFind() Find[T]            { return di.Single(find.New[T], d) }
func (d *Deps[T]) TargetFinder() Finder[T]        { return portals.NewFind[T] }
func (d *Deps[T]) NewTray() NewTray               { return tray.NewRun }
func (d *Deps[T]) RunSpawn() Dispatch             { return di.Single(spawn.Inject[T], d) }
func (d *Deps[T]) JoinTarget() Dispatch           { return query.NewRunner[Portal](PortOpen).Run }
func (d *Deps[T]) FeatDispatch() Dispatch         { return di.Single(dispatch.Inject, d) }
func (d *Deps[T]) DispatchService() Dispatch      { return di.Single(serve.Inject[T], d).Dispatch }
func (d *Deps[T]) Astral() serve.Astral           { return serve.CheckAstral }
func (d *Deps[T]) RunService() serve.Service      { return service.NewRun }
func (d *Deps[T]) FeatObserve() serve.Observe     { return featApps.Observe }
func (d *Deps[T]) FeatCreate() *create.Feat {
	return create.NewFeat(template.NewRun, d.FeatBuild().Dist)
}
func (d *Deps[T]) FeatBuild() *build.Feat {
	return build.NewFeat(
		dist.NewRun, pack.Run,
		sources.FromFS[NodeModule](js.PortalLibFS),
	)
}
func (d *Deps[T]) RpcHandlers() rpc.Handlers {
	return rpc.Handlers{
		PortMsg.Name: msg.NewBroadcast(PortMsg, d.Processes()).BroadcastMsg,
	}
}
func (d *Deps[T]) TargetRun() Run[T] {
	return multi.NewRunner[T](
		app.Run(exec.NewPortal[PortalJs]("portal-dev-goja", "o").Run),
		app.Run(exec.NewPortal[PortalHtml]("portal-dev-wails", "o").Run),
		app.Run(exec.NewPortal[ProjectGo]("portal-dev-go", "o").Run),
		app.Run(exec.NewPortal[AppExec]("portal-dev-exec", "o").Run),
	).Run
}
