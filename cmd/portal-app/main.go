package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/clir"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/runner/tray"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/cache"
	"github.com/cryptopunkscc/portal/target/find"
	"github.com/cryptopunkscc/portal/target/spawn"
	"sync"
)

func main() {
	deps := NewDeps[App]()
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app").Set(&ctx)
	cli := clir.NewCli(ctx,
		"Portal-app",
		"Portal applications service.",
		version.Run,
	)
	cli.Serve(deps.FeatServe())
	cli.Apps(deps.TargetFind())
	cli.List(featApps.List)
	cli.Install(featApps.Install)
	cli.Uninstall(featApps.Uninstall)
	go singal.OnShutdown(cancel)
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	deps.WaitGroup().Wait()
}

type Deps[T App] struct {
	di.Cache
	wg        sync.WaitGroup
	processes sig.Map[string, T]
	cache     *Cache[T]
}

func NewDeps[T App]() *Deps[T]                    { return &Deps[T]{cache: NewCache[T]()} }
func (d *Deps[T]) Executable() string             { return "portal" }
func (d *Deps[T]) GetCacheDir() string            { return di.Single(cache.Dir, d) }
func (d *Deps[T]) WaitGroup() *sync.WaitGroup     { return &d.wg }
func (d *Deps[T]) Processes() *sig.Map[string, T] { return &d.processes }
func (d *Deps[T]) Port() Port                     { return PortPortal }
func (d *Deps[T]) Path() Path                     { return featApps.Path }
func (d *Deps[T]) NewTray() NewTray               { return tray.NewRun }
func (d *Deps[T]) RunSpawn() Dispatch             { return di.Single(spawn.Inject[T], d) }
func (d *Deps[T]) TargetFind() Find[T]            { return di.Single(find.New[T], d) }
func (d *Deps[T]) TargetFinder() Finder[T]        { return apps.NewFind[T] }
func (d *Deps[T]) TargetCache() *Cache[T]         { return d.cache }
func (d *Deps[T]) Astral() serve.Astral           { return exec.Astral }
func (d *Deps[T]) RunService() serve.Service      { return service.NewRun }
func (d *Deps[T]) FeatObserve() serve.Observe     { return featApps.Observe }
func (d *Deps[T]) FeatServe() clir.Serve          { return di.Single(serve.Inject[T], d).Run }
func (d *Deps[T]) RpcHandlers() rpc.Handlers      { return nil }
func (d *Deps[T]) TargetRun() Run[T] {
	return multi.NewRunner[T](
		app.Run(exec.NewPortal[AppJs]("portal-app-goja", "o").Run),
		app.Run(exec.NewPortal[AppHtml]("portal-app-wails", "o").Run),
		app.Run(exec.NewBundleRunner(d.GetCacheDir()).Run),
	).Run
}
