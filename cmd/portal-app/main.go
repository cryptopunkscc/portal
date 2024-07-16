package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/clir"
	feature "github.com/cryptopunkscc/portal/feat"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/runner/tray"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app").Set(&ctx)

	go singal.OnShutdown(cancel)

	scope := feature.Scope[target.App]{
		Executable:    "portal",
		Port:          target.PortPortal,
		Astral:        exec.Astral,
		GetPath:       featApps.Path,
		FeatObserve:   featApps.Observe,
		TargetFinder:  apps.NewFind[target.App],
		TargetCache:   target.NewCache[target.App](),
		Processes:     &sig.Map[string, target.App]{},
		NewRunTray:    tray.NewRun,
		NewRunService: service.NewRun,
		NewExecTarget: func(_ string, cache string) target.Run[target.App] {
			return multi.NewRunner[target.App](
				app.Run(exec.NewPortal[target.AppJs]("portal-app-goja", "o").Run),
				app.Run(exec.NewPortal[target.AppHtml]("portal-app-wails", "o").Run),
				app.Run(exec.NewBundleRunner(cache).Run),
			).Run
		},
	}

	cli := clir.NewCli(ctx,
		"Portal-app",
		"Portal applications service.",
		version.Run,
	)
	cli.Serve(scope.GetServeFeature().Run)
	cli.Apps(scope.GetTargetFind())
	cli.List(featApps.List)
	cli.Install(featApps.Install)
	cli.Uninstall(featApps.Uninstall)

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	scope.WaitGroup.Wait()
}
