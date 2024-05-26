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
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
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
	"os"
	"sync"
)

func main() {
	println("...")
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().I().Set(&ctx)
	log.Scope("main").Println("starting portal", os.Args)
	go osexec.OnShutdown(cancel)

	wait := &sync.WaitGroup{}
	executable := "portal"
	port := "portal"
	portOpen := "portal.open"
	findApps := createAppsFind()
	runQuery := query.NewRunner[target.App](portOpen).Run

	featDispatch := dispatch.NewFeat(executable, runQuery)
	featServe := createServeFeature(wait, executable, port, findApps)
	featOpen := createOpenFeature(runQuery, findApps)

	cli := clir.NewCli(ctx, manifest.Name, manifest.Description, version.Run)

	cli.Dispatch(featDispatch)
	cli.Serve(featServe)
	cli.Open(featOpen)
	cli.List(featApps.List)
	cli.Install(featApps.Install)
	cli.Uninstall(featApps.Uninstall)
	cli.Apps(findApps)

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
	queryOpen target.Dispatch,
	findApps target.Find[target.App],
) target.Dispatch {
	newApphost := apphost.NewFactory(queryOpen)
	newApi := target.ApiFactory(NewAdapter,
		newApphost.NewAdapter,
		newApphost.WithTimeout,
	)
	runApp := app.NewRun(newApi)

	return open.NewFeat[target.App](findApps, runApp)
}

func createServeFeature(
	wait *sync.WaitGroup,
	executable string,
	port string,
	findApps target.Find[target.App],
) func(context.Context, bool) error {
	runProc := exec.NewRun[target.App](executable)
	runSpawn := spawn.NewRunner(wait, findApps, runProc).Run
	return serve.NewFeat(port, runSpawn, tray.New(runSpawn))
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

	return target.Cached(apps.NewFind)(findPath, embedApps.LauncherSvelteFS)
}
