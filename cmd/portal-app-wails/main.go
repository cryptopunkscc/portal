package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/app"
	factory "github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	mod := &Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-wails").Set(&ctx)
	go sig.OnShutdown(log, cancel)

	err := cli.New(cmd.Handler{
		Name: "portal-app-wails",
		Desc: "Portal html runner driven by wails.",
		Sub: cmd.Handlers{
			{
				Func: open.Runner[AppHtml](mod),
				Name: "o",
				Desc: "Start portal app in wails runner.",
				Params: cmd.Params{
					{Type: "string", Desc: "Absolute path to app bundle or directory."},
				},
			},
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}).Run(ctx)

	if err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ app.Module[AppHtml] }
type Adapter struct{ bind.Runtime }

func (d *Module) Runner() Run[AppHtml] { return wails.NewRun(d.runtime) }
func (d *Module) runtime(ctx context.Context, portal Portal_) bind.Runtime {
	return &Adapter{factory.FrontendRuntime()(ctx, portal)}
}
