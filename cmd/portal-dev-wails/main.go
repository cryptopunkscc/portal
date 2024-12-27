package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	factory "github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/factory/dev"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails_dist"
	"github.com/cryptopunkscc/portal/runner/wails_pro"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	mod := &Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-wails").Set(&ctx)
	go sig.OnShutdown(cancel)

	err := cli.New(cmd.Handler{
		Name: "portal-dev-wails",
		Desc: "Portal html development driven by wails.",
		Sub: cmd.Handlers{
			{
				Func: open.Feat[PortalHtml](mod),
				Name: "o",
				Desc: "Start portal js app development.",
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

type Module struct{ dev.Module[PortalHtml] }
type Adapter struct{ bind.Runtime }

func (d *Module) Runner() Run[PortalHtml] {
	return multi.Runner[PortalHtml](
		reload.Immutable(d.runtime, PortMsg, wails_pro.NewRunner), // FIXME propagate sendMsg
		reload.Mutable(d.runtime, PortMsg, wails_dist.NewRunner),
		reload.Immutable(d.runtime, PortMsg, wails.NewRunner),
	)
}
func (d *Module) runtime(ctx context.Context, portal Portal_) bind.Runtime {
	a := &Adapter{factory.FrontendRuntime()(ctx, portal)}
	return a
}
