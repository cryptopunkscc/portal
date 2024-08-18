package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/app"
	factory "github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runtime/bind"
)

func main() {
	m := &Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-wails").Set(&ctx)
	go sig.OnShutdown(cancel)
	cli := clir.NewCli(ctx,
		"Portal-wails",
		"Portal html runner driven by wails.",
		version.Run,
	)
	cli.Open(open.Feat[AppHtml](m))

	if err := cli.Run(); err != nil {
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
