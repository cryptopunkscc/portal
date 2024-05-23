package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	embedApps "github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/clir"
	featApps "github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/dispatch"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	"github.com/cryptopunkscc/go-astral-js/mock/appstore"
	osexec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/app"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/runner/spawn"
	"github.com/cryptopunkscc/go-astral-js/runner/tray"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apphost"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/portal"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"log"
	"os"
	"sync"
)

func main() {
	log.Println("starting portal", os.Args)
	ctx, cancel := context.WithCancel(context.Background())
	go osexec.OnShutdown(cancel)

	wait := &sync.WaitGroup{}
	executable := "portal"
	findApps := createAppsFind()

	featDispatch := dispatch.NewFeat(executable)
	featServe := createServeFeature(wait, executable, findApps)
	featOpen := createOpenFeature(findApps)

	cli := clir.NewCli(ctx, manifest.Name, manifest.Description, version.Run)

	cli.Dispatch(featDispatch)
	cli.Serve(featServe)
	cli.Open(featOpen)
	cli.List(featApps.List)
	cli.Install(featApps.Install)
	cli.Uninstall(featApps.Uninstall)

	err := cli.Run()
	cancel()
	if err != nil {
		log.Println(err)
	}
	wait.Wait()
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }

func createOpenFeature(
	findApps target.Find[target.App],
) target.Dispatch {
	runQuery := query.NewRunner[target.App]().Run
	newApphost := apphost.NewFactory(runQuery)
	newApi := target.ApiFactory(NewAdapter,
		newApphost.NewAdapter,
		newApphost.WithTimeout,
	)
	runApp := app.NewRunner(newApi)

	return open.NewFeat[target.App](findApps, runApp)
}

func createServeFeature(
	wait *sync.WaitGroup,
	executable string,
	findApps target.Find[target.App],
) func(context.Context, bool) error {
	runProc := exec.NewRunner[target.App](executable)
	runSpawn := spawn.NewRunner(wait, findApps, runProc).Run
	return serve.NewFeat(runSpawn, tray.NewRunner(runSpawn))
}

func createAppsFind() target.Find[target.App] {
	resolveEmbed := portal.NewResolver[target.App](
		apps.Resolve[target.App](),
		source.FromFS(embedApps.LauncherSvelteFS),
	)
	findPath := target.Mapper[string, string](
		resolveEmbed.Path,
		appstore.Path,
	)
	return apps.NewFinder(findPath, embedApps.LauncherSvelteFS).Find
}
