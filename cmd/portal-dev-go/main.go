package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/factory/dev"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/go_dev"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	mod := &Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-go").Set(&ctx)
	go sig.OnShutdown(cancel)

	err := cli.New(cmd.Handler{
		Name: "portal-dev-go",
		Desc: "Portal go development.",
		Sub: cmd.Handlers{
			{
				Func: open.Runner[ProjectGo](mod),
				Name: "o",
				Desc: "Start portal golang app development.",
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

type Module struct{ dev.Module[ProjectGo] }

func (d *Module) Runner() Run[ProjectGo] {
	return multi.Runner[ProjectGo](
		reload.Mutable(bind.DefaultRuntime(), PortMsg, go_dev.Adapter(exec.Dist().Run)),
	)
}
