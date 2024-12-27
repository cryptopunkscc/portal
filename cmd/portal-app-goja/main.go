package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/app"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	mod := &Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-goja").Set(&ctx)
	go sig.OnShutdown(cancel)

	err := cli.New(cmd.Handler{
		Name: "portal-app-goja",
		Desc: "Portal js runner driven by goja.",
		Sub: cmd.Handlers{
			{
				Func: open.Feat[AppJs](mod),
				Name: "o",
				Desc: "Start portal app in goja runner.",
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

type Module struct{ app.Module[AppJs] }

func (d *Module) Runner() Run[AppJs] { return goja.NewRun(bind.BackendRuntime()) }
