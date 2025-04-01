package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails_dist"
	"github.com/cryptopunkscc/portal/runner/wails_pro"
)

func main() { cli.Run(Application[PortalHtml]{}.handler()) }

type Application[T PortalHtml] struct{}

func (a Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Func: open.NewRun[T](&a),
		Name: "dev-html",
		Desc: "Start html app development in wails runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}

func (a Application[T]) Runner() Run[T] {
	return multi.Runner[T](
		reload.Immutable(a.core, wails_pro.ReRunner), // FIXME propagate sendMsg
		reload.Mutable(a.core, wails_dist.ReRunner),
		reload.Immutable(a.core, wails.ReRunner),
	)
}
func (a Application[T]) Resolver() Resolve[T] { return sources.Resolver[T]() }

func (a Application[T]) core(ctx context.Context, portal Portal_) (bind.Core, context.Context) {
	r, ctx := bind.NewFrontendCoreFunc()(ctx, portal)
	return &Adapter{r}, ctx
}

type Adapter struct{ bind.Core }
