package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/core/rpc/cmd"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runner/wails"
)

func main() { cli.Run(Application[AppHtml]{}.Handler()) }

type Application[T AppHtml] struct{}

func (a Application[T]) Handler() cmd.Handler {
	return cmd.Handler{
		Func: open.Runner[T](&a),
		Name: "html",
		Desc: "Start portal app in wails runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}
func (a Application[T]) Resolver() Resolve[T] { return apps.Resolver[T]() }
func (a Application[T]) Runner() Run[T]       { return multi.Runner[T](app.Runner(wails.Runner(a.runtime))) }
func (a Application[T]) runtime(ctx context.Context, portal Portal_) (bind.Runtime, context.Context) {
	r, ctx := bind.FrontendRuntime()(ctx, portal)
	return &Adapter{r}, ctx
}

type Adapter struct{ bind.Runtime }
