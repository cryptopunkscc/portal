package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/di/run/dev"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/reload"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/cache"
)

func main() {
	mod := Module{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-exec").Set(&ctx)
	go sig.OnShutdown(cancel)
	cli := clir.NewCli(ctx,
		"Portal-dev-exec",
		"Portal js development runner for executables.",
		version.Run,
	)
	cli.Open(mod.FeatOpen())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Module struct{ dev.Module[AppExec] }

func (d *Module) Executable() string  { return "portal-dev" }
func (d *Module) GetCacheDir() string { return di.Single(cache.Dir, cache.Deps(d)) }
func (d *Module) WrapApi(api Api) Api { return api }
func (d *Module) NewRunTarget(newApi NewApi) Run[AppExec] {
	return multi.NewRunner[AppExec](
		reload.Immutable(newApi, PortMsg, reload.Adapter(exec.NewBundleRunner(d.GetCacheDir()))),
		reload.Immutable(newApi, PortMsg, reload.Adapter(exec.NewDistRunner())),
	).Run
}
