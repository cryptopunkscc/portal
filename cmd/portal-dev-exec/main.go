package main

import (
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	osExec "github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"golang.org/x/net/context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	log := plog.New().D().Scope("dev-exec").Set(&ctx)

	scope := feature.Scope[target.AppExec]{
		WrapApi:        NewAdapter,
		GetPath:        featApps.Path,
		TargetFinder:   apps.NewFind[target.AppExec],
		TargetCache:    target.NewCache[target.AppExec](),
		DispatchTarget: query.NewRunner[target.AppExec](target.PortOpen).Start,
	}
	scope.NewRunTarget = func(newApi target.NewApi) target.Run[target.AppExec] {
		return multi.NewRunner[target.AppExec](
			reload.Immutable(newApi, target.PortMsg, reload.Adapter(exec.NewBundleRunner(scope.GetCacheDir()))),
			reload.Immutable(newApi, target.PortMsg, reload.Adapter(exec.NewDistRunner())),
		).Run
	}
	cli := clir.NewCli(ctx,
		"Portal-dev-exec",
		"Portal js development runner for executables.",
		version.Run,
	)
	cli.Open(scope.GetOpenFeature())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }
