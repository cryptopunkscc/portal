package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/clir"
	featApps "github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/dispatch"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	osexec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/app"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/runner/spawn"
	"github.com/cryptopunkscc/go-astral-js/runner/tray"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apphost"
	apps "github.com/cryptopunkscc/go-astral-js/target/apps"
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

	findApps := apps.Find(featApps.Path)

	runQuery := query.NewRunner[target.App]().Run
	newApphost := apphost.NewFactory(runQuery)
	newApi := target.ApiFactory(NewAdapter,
		newApphost.NewAdapter,
		newApphost.WithTimeout,
	)
	runApp := app.NewRunner(newApi)
	runProc := exec.NewRunner[target.App](executable)
	runSpawn := spawn.NewRunner(wait, findApps, runProc).Run

	featDispatch := dispatch.NewFeat(executable)
	featServe := serve.NewFeat(runSpawn, tray.NewRunner(runSpawn))
	featOpen := open.NewFeat[target.App](findApps, runApp)

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
