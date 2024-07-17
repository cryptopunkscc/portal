package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/di/run/app"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/goja"
	. "github.com/cryptopunkscc/portal/target"
)

func main() {
	mod := Module{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-goja").Set(&ctx)
	go sig.OnShutdown(cancel)
	cli := clir.NewCli(ctx,
		"Portal-goja",
		"Portal js runner driven by goja.",
		version.Run,
	)
	cli.Open(mod.FeatOpen())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ app.Module[AppJs] }

func (d *Module) WrapApi(api Api) Api                   { return api }
func (d *Module) NewRunTarget(newApi NewApi) Run[AppJs] { return goja.NewRun(newApi) }
