package main

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/portal"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/unknown"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/find"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/serve"
	"github.com/cryptopunkscc/portal/runner/supervisor"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runtime/apps"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	application := Application[Portal_]{}
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

type Application[T Portal_] struct {
	CancelFunc context.CancelFunc
	wg         sync.WaitGroup
	processes  sig.Map[string, T]
	cache      Cache[T]
}

func (a *Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Name: "portald",
		Desc: "Start portal applications service",
		Func: serve.Runner(a),
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}
}

func (a *Application[T]) Open() Run[portal.OpenOpt] {
	return func(ctx context.Context, opt portal.OpenOpt, cmd ...string) (err error) {
		if len(cmd) == 0 {
			return errors.New("no command")
		}
		src := cmd[0]
		args := cmd[1:]

		var schemaPrefix []string
		if opt.Schema != "" {
			schemaPrefix = []string{opt.Schema}
		}
		plog.Get(ctx).Type(a).Println("open:", opt, cmd, opt.Order)
		return find.Runner[T](
			FindByPath(
				source.File, Any[T](
					Skip("node_modules"),
					Try(exec2.ResolveBundle),
					Try(exec2.ResolveDist),
					Try(unknown.ResolveBundle),
					Try(unknown.ResolveDist),
					Try(unknown.ResolveProject),
				),
			).ById(path.Resolver(apps.Source)).
				Cached(&a.cache).
				Reduced(a.priority(opt.Order)...),
			supervisor.Runner[T](
				&a.wg,
				&a.processes,
				multi.Runner[T](
					app.Runner(exec.BundleRunner(a.cacheDir())),
					app.Runner(exec.DistRunner()),
					app.Runner(exec.AnyRunner(a.cacheDir(), schemaPrefix...)),
				),
			),
		).Call(ctx, src, args...)
	}
}

func (a *Application[T]) Shutdown() context.CancelFunc        { return a.CancelFunc }
func (a *Application[T]) Astral() serve.Astral                { return exec.Astral }
func (a *Application[T]) cacheDir() string                    { return CacheDir("portal") }
func (a *Application[T]) priority(order []int) (out Priority) { return matchers.Sort(order) }

var matchers = Priority{
	Match[Bundle_],
	Match[Dist_],
	Match[Project_],
}
