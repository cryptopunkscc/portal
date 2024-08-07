package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/di/run/app"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/wails"
	. "github.com/cryptopunkscc/portal/target"
)

func main() {
	module := Module{}
	module.Deps = &module
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-wails").Set(&ctx)
	go sig.OnShutdown(cancel)
	cli := clir.NewCli(ctx,
		"Portal-wails",
		"Portal html runner driven by wails.",
		version.Run,
	)
	cli.Open(module.FeatOpen())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ app.Module[AppHtml] }
type Adapter struct{ Api }

func (d *Module) WrapApi(api Api) Api                     { return &Adapter{api} }
func (d *Module) NewRunTarget(newApi NewApi) Run[AppHtml] { return wails.NewRun(newApi) }
