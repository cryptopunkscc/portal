package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/factory/dev"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/reload"
	. "github.com/cryptopunkscc/portal/target"
)

func main() {
	mod := &Module{}
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-exec").Set(&ctx)
	go sig.OnShutdown(cancel)
	cli := clir.NewCli(ctx,
		"Portal-dev-exec",
		"Portal js development runner for executables.",
		version.Run,
	)
	cli.Open(open.Feat[AppExec](mod))
	if err := cli.Run(); err != nil {
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
