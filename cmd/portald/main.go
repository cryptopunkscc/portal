package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/portald"
	"github.com/cryptopunkscc/portal/runner/version"
	_ "github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	application := Application[Portal_]{}
	application.Shutdown = cancel
	application.CacheDir = CacheDir("portal")
	log := plog.New().D().Scope("portald").Set(&ctx)
	go singal.OnShutdown(log, cancel)
	handler := application.handler()
	cmd.InjectHelp(&handler)
	err := cli.New(handler).Run(ctx)
	if err != nil {
		log.E().Println("finished with error:", err)
	}
}

type Application[T Portal_] struct {
	portald.Runner[Portal_]
}

func (a *Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Name: "portald",
		Desc: "Portal daemon.",
		Func: a.run,
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}

func (a *Application[T]) run(ctx context.Context) (err error) {
	if err = exec.Astral(ctx); err != nil {
		return
	}
	err = a.Run(ctx)
	return
}
