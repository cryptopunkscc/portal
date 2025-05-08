package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/target/source"
)

func main() { cli.Run(Application{}.Handler()) }

type Application struct{}

func (a Application) Handler() cmd.Handler {
	return cmd.Handler{
		Func: source.File.NewRun(wails.Runner(a.core).Try),
		Name: "html",
		Desc: "Start portal app in wails runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Name},
		},
	}
}

func (a Application) core(ctx context.Context, portal Portal_) (bind.Core, context.Context) {
	r, ctx := bind.NewFrontendCore(ctx, portal)
	return &Adapter{r}, ctx
}

type Adapter struct{ bind.Core }
