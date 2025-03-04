package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	factory "github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails_dist"
	"github.com/cryptopunkscc/portal/runner/wails_pro"
	_ "github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
)

func main() { cli.Run(Application[PortalHtml]{}.handler()) }

type Application[T PortalHtml] struct{}

func (a Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Func: open.Runner[T](&a),
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
		reload.Immutable(a.runtime, wails_pro.ReRunner), // FIXME propagate sendMsg
		reload.Mutable(a.runtime, wails_dist.ReRunner),
		reload.Immutable(a.runtime, wails.ReRunner),
	)
}
func (a Application[T]) Resolver() Resolve[T] { return sources.Resolver[T]() }

func (a Application[T]) runtime(ctx context.Context, portal Portal_) (bind.Runtime, context.Context) {
	r, ctx := factory.FrontendRuntime()(ctx, portal)
	return &Adapter{r}, ctx
}

type Adapter struct{ bind.Runtime }
