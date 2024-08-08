package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/srv"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	. "github.com/cryptopunkscc/portal/target"
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
	cli.Serve(mod.Serve())
	go singal.OnShutdown(cancel)
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	mod.WaitGroup().Wait()
}

type Module[T App_] struct{ srv.Module[T] }

func (d *Module[T]) Executable() string   { return "portal" }
func (d *Module[T]) Serve() clir.Serve    { return serve.Run(d) }
func (d *Module[T]) Astral() serve.Astral { return exec.Astral }
func (d *Module[T]) Resolve() Resolve[T]  { return apps.Resolver[T]() }
func (d *Module[T]) Run() Run[T] {
	return multi.NewRunner[T](
		app.Run(exec.NewPortal[AppJs]("portal-app-goja", "o").Run),
		app.Run(exec.NewPortal[AppHtml]("portal-app-wails", "o").Run),
		app.Run(exec.NewBundleRunner(CacheDir(d.Executable())).Run),
	).Run
}
func (d *Module[T]) Priority() Priority {
	return []Matcher{
		Match[Bundle_],
		Match[Dist_],
	}
}
