package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/di/srv"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/cache"
)

func main() {
	mod := Module[App_]{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	mod.CancelFunc = cancel
	log := plog.New().D().Scope("app").Set(&ctx)
	cli := clir.NewCli(ctx,
		"Portal-app",
		"Portal applications service.",
		version.Run,
	)
	cli.Serve(mod.FeatServe())
	cli.List(featApps.List)
	cli.Install(featApps.Install)
	cli.Uninstall(featApps.Uninstall)
	go singal.OnShutdown(cancel)
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	mod.WaitGroup().Wait()
}

type Module[T App_] struct{ srv.Module[App_] }

func (d *Module[T]) Executable() string        { return "portal" }
func (d *Module[T]) Astral() serve.Astral      { return exec.Astral }
func (d *Module[T]) RpcHandlers() rpc.Handlers { return nil }
func (d *Module[T]) TargetResolve() Resolve[T] { return apps.Resolver[T]() }
func (d *Module[T]) TargetRun() Run[T] {
	return multi.NewRunner[T](
		app.Run(exec.NewPortal[AppJs]("portal-app-goja", "o").Run),
		app.Run(exec.NewPortal[AppHtml]("portal-app-wails", "o").Run),
		app.Run(exec.NewBundleRunner(d.CacheDir()).Run),
	).Run
}
func (d *Module[T]) Priority() Priority {
	return []Matcher{
		Match[Bundle_],
		Match[Dist_],
	}
}
func (d *Module[T]) CacheDir() string      { return di.S(cache.Dir, cache.Deps(d)) }
func (d *Module[T]) FeatServe() clir.Serve { return serve.Inject[T](d).Run }
