package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/app"
	factory "github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() { cli.Run(Application{}.Handler()) }

type Application struct{ app.Module[AppHtml] }

func (a Application) Handler() cmd.Handler {
	return cmd.Handler{
		Func: open.Runner[AppHtml](&a),
		Name: "html",
		Desc: "Start portal app in wails runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}
}

func (a Application) Runner() Run[AppHtml] { return wails.NewRun(a.runtime) }
func (a Application) runtime(ctx context.Context, portal Portal_) bind.Runtime {
	return &Adapter{factory.FrontendRuntime()(ctx, portal)}
}

type Adapter struct{ bind.Runtime }
