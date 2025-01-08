package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/factory/dev"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/exec"
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
	log := plog.New().D().Scope("dev-exec").Set(&ctx)
	go sig.OnShutdown(log, cancel)

	err := cli.New(cmd.Handler{
		Name: "portal-dev-exec",
		Desc: "Portal development runner for executables.",
		Sub: cmd.Handlers{
			{
				Func: open.Runner[AppExec](mod),
				Name: "o",
				Desc: "Start portal app executable.",
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

type Module struct{ dev.Module[AppExec] }

func (d *Module) Runner() Run[AppExec] {
	return multi.Runner[AppExec](
		reload.Immutable(bind.DefaultRuntime(), PortMsg, reload.Adapter(exec.Bundle(CacheDir("portal-dev")))),
		reload.Immutable(bind.DefaultRuntime(), PortMsg, reload.Adapter(exec.Dist())),
	)
}
