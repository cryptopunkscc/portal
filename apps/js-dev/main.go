package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/goja_pro"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() { cli.Run(Application[PortalJs]{}.handler()) }

type Application[T PortalJs] struct{}

func (a Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Func: open.Runner[T](&a),
		Name: "dev-js",
		Desc: "Start portal js app development in goja runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}
}

func (a Application[T]) Runner() Run[T] {
	return multi.Runner[T](
		reload.Mutable(bind.BackendRuntime(), goja_pro.NewRunner),
		reload.Mutable(bind.BackendRuntime(), goja_dist.NewRunner),
		reload.Immutable(bind.BackendRuntime(), goja.NewRunner),
	)
}

func (a Application[T]) Resolver() Resolve[T] { return sources.Resolver[T]() }
