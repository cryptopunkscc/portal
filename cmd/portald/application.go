package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/portald"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/astrald"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/version"
)

type Application[T Portal_] struct {
	portald.Runner[Portal_]
}

func (a *Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Func: a.run,
		Name: "portald",
		Desc: "Portal daemon.",
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}

func (a *Application[T]) run(ctx context.Context) (err error) {
	runner := astrald.Runner{Astrald: &exec.Astrald{}}
	if err = runner.Start(ctx); err != nil {
		return
	}
	err = a.Run(ctx)
	return
}
