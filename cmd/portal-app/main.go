package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/srv"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	mod := Module[App_]{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	mod.CancelFunc = cancel
	log := plog.New().D().Scope("app").Set(&ctx)
	go singal.OnShutdown(cancel)

	err := cli.New(cmd.Handler{
		Name: "portal-app",
		Desc: "Portal applications service",
		Sub: cmd.Handlers{
			{Name: "serve s", Desc: "Start portal daemon", Func: mod.Serve()},
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}).Run(ctx)

	if err != nil {
		log.Println(err)
	}
	cancel()
	mod.WaitGroup().Wait()
}

type Module[T App_] struct{ srv.Module[T] }

func (d *Module[T]) Executable() string   { return "portal" }
func (d *Module[T]) Serve() Request       { return serve.Feat(d) }
func (d *Module[T]) Astral() serve.Astral { return exec.Astral }
func (d *Module[T]) Resolve() Resolve[T]  { return apps.Resolver[T]() }
func (d *Module[T]) Run() Run[T] {
	return multi.Runner[T](
		app.Run(exec.Portal[AppJs]("portal-app-goja", "o").Run),
		app.Run(exec.Portal[AppHtml]("portal-app-wails", "o").Run),
		app.Run(exec.Bundle(CacheDir(d.Executable())).Run),
	)
}
func (d *Module[T]) Priority() Priority {
	return []Matcher{
		Match[Bundle_],
		Match[Dist_],
	}
}
