package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/srv"
	"github.com/cryptopunkscc/portal/pkg/plog"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/unknown"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/serve"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	application := portald[App_]{}
	application.Deps = &application
	application.CancelFunc = cancel
	log := plog.New().D().Scope("portald").Set(&ctx)
	go singal.OnShutdown(cancel)

	err := cli.New(application.Handler()).Run(ctx)
	if err != nil {
		log.Println(err)
	}
	cancel()
	application.WaitGroup().Wait()
}

type portald[T App_] struct{ srv.Module[T] }

func (d *portald[T]) Handler() cmd.Handler {
	return cmd.Handler{
		Name: "portald",
		Desc: "Start portal applications service",
		Func: serve.Runner(d),
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}
}

func (d *portald[T]) Astral() serve.Astral { return exec.Astral }

func (d *portald[T]) Run() Run[T] {
	return multi.Runner[T](
		app.Run(exec.Bundle(CacheDir("portal")).Run),
		app.Run(exec.Dist().Run),
		app.Run(exec.Any(d.runner).Run),
	)
}
func (d *portald[T]) Priority() Priority {
	return []Matcher{
		Match[Bundle_],
		Match[Dist_],
	}
}
func (d *portald[T]) Resolve() Resolve[T] {
	return Any[T](
		Skip("node_modules"),
		Try(exec2.ResolveBundle),
		Try(exec2.ResolveDist),
		Try(unknown.ResolveBundle),
		Try(unknown.ResolveDist),
	)
}

// TODO resolve dynamically
func (d *portald[T]) runner(script string) string {
	return map[string]string{
		"js":   "portal-app-goja",
		"html": "portal-app-wails",
	}[script]
}
