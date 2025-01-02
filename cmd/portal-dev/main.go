package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/srv"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/start"
	"github.com/cryptopunkscc/portal/pkg/plog"
	signal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/request/query"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/resolve/unknown"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runtime/msg"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	mod := deps[Portal_]{}
	mod.Deps = &mod
	mod.CancelFunc = cancel
	go signal.OnShutdown(cancel)
	println("...")
	plog.ErrorStackTrace = true
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")

	err := cli.New(cmd.Handler{
		Name: "portal-dev",
		Desc: "Start portal project or app from a given source in development environment.",
		Params: cmd.Params{
			{Type: "string", Desc: "Application source. The source can be a app name, package name, app bundle path or app dir."},
		},
		Func: start.Feat(&mod),
	}).Run(ctx)

	if err != nil {
		log.Println(err)
	}
	cancel()
	mod.WaitGroup().Wait()
}

type deps[T Portal_] struct{ srv.Module[T] }

func (d *deps[T]) Executable() string   { return "portal-dev" }
func (d *deps[T]) Request() Request     { return query.Request.Run }
func (d *deps[T]) Serve() Request       { return serve.Feat(d).Start }
func (d *deps[T]) Astral() serve.Astral { return serve.CheckAstral }
func (d *deps[T]) Handlers() cmd.Handlers {
	return cmd.Handlers{
		{Name: PortMsg.Name(), Func: msg.NewBroadcast(PortMsg, d.Processes()).BroadcastMsg},
	}
}
func (d *deps[T]) Run() Run[T] {
	return multi.Runner[T](
		app.Run(exec.Any(d.runner).Run),
	)
}
func (d *deps[T]) Priority() Priority {
	return Priority{
		Match[Project_],
		Match[Dist_],
		Match[Bundle_],
	}
}
func (d *deps[T]) Resolve2() Resolve[T] { return sources.Resolver[T]() }

func (d *deps[T]) Resolve() Resolve[T] {
	return Any[T](
		Skip("node_modules"),
		Try(exec2.ResolveDist),
		Try(exec2.ResolveBundle),
		Try(unknown.ResolveDist),
		Try(unknown.ResolveBundle),
		Try(unknown.ResolveProject),
	)
}

// TODO resolve dynamically
func (d *deps[T]) runner(script string) string {
	return map[string]string{
		"js":   "portal-dev-goja",
		"html": "portal-dev-wails",
		"go":   "portal-dev-go",
		"exec": "portal-dev-exec",
	}[script]
}
