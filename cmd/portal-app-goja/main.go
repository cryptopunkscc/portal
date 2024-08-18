package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/app"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/goja"
)

func main() {
	mod := &Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-goja").Set(&ctx)
	go sig.OnShutdown(cancel)
	cli := clir.NewCli(ctx,
		"Portal-goja",
		"Portal js runner driven by goja.",
		version.Run,
	)
	cli.Open(open.Feat[AppJs](mod))
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ app.Module[AppJs] }

func (d *Module) Runner() Run[AppJs] { return goja.NewRun(bind.BackendRuntime()) }
