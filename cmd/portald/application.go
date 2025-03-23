package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/portald"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/version"
)

type Application[T Portal_] struct {
	portald.Service[Portal_]
}

func (a *Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Func: a.Run,
		Name: "portald",
		Desc: "Portal daemon.",
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}
