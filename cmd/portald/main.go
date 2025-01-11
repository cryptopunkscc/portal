package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/apphost"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/pkg/plog"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/unknown"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/find"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/serve"
	"github.com/cryptopunkscc/portal/runner/supervisor"
	"github.com/cryptopunkscc/portal/runner/version"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	application := Application[App_]{}
	application.CancelFunc = cancel
	log := plog.New().D().Scope("portald").Set(&ctx)
	go singal.OnShutdown(log, cancel)

	err := cli.New(application.handler()).Run(ctx)
	if err != nil {
		log.Println(err)
	}
	cancel()
	application.wg.Wait()
}

type Application[T App_] struct {
	CancelFunc context.CancelFunc
	wg         sync.WaitGroup
	processes  sig.Map[string, T]
	cache      Cache[T]
}

func (d *Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Name: "portald",
		Desc: "Start portal applications service",
		Func: serve.Runner(d),
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}
}

func (d *Application[T]) Open() Run[string] {
	return find.Runner[T](
		FindByPath(
			source.File, Any[T](
				Skip("node_modules"),
				Try(exec2.ResolveBundle),
				Try(exec2.ResolveDist),
				Try(unknown.ResolveBundle),
				Try(unknown.ResolveDist),
			),
		).ById(appstore.Path).Cached(&d.cache).Reduced(
			Match[Bundle_],
			Match[Dist_],
		),
		supervisor.Runner[T](
			&d.wg,
			&d.processes,
			multi.Runner[T](
				app.Run(exec.BundleRun(CacheDir("portal"))),
				app.Run(exec.Dist().Run),
				app.Run(exec.AnyRun(CacheDir("portal"))),
			),
		),
	)
}

func (d *Application[T]) Shutdown() context.CancelFunc                   { return d.CancelFunc }
func (d *Application[T]) Observe() func(context.Context, rpc.Conn) error { return appstore.Observe }

func (d *Application[T]) Port() apphost.Port     { return PortPortal }
func (d *Application[T]) Astral() serve.Astral   { return exec.Astral }
func (d *Application[T]) Handlers() cmd.Handlers { return cmd.Handlers{} }
