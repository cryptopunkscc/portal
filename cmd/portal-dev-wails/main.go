package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/run/dev"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails_dist"
	"github.com/cryptopunkscc/portal/runner/wails_pro"
	. "github.com/cryptopunkscc/portal/target"
)

func main() {
	mod := Module{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-wails").Set(&ctx)
	go sig.OnShutdown(cancel)
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx,
		"Portal-dev-wails",
		"Portal html development driven by wails.",
		version.Run,
	)
	cli.Open(mod.FeatOpen())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ dev.Module[PortalHtml] }
type Adapter struct{ Api }

func NewAdapter(api Api) Api          { return &Adapter{Api: api} }
func (d *Module) WrapApi(api Api) Api { return NewAdapter(api) }
func (d *Module) NewRunTarget(newApi NewApi) Run[PortalHtml] {
	return multi.NewRunner[PortalHtml](
		reload.Immutable(newApi, PortMsg, wails_pro.NewRunner), // FIXME propagate sendMsg
		reload.Mutable(newApi, PortMsg, wails_dist.NewRunner),
		reload.Immutable(newApi, PortMsg, wails.NewRunner),
	).Run
}
