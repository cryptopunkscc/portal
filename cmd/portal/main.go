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
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/apps"
	"github.com/cryptopunkscc/go-astral-js/runner/app"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/spawn"
	"github.com/cryptopunkscc/go-astral-js/runner/tray"
	"log"
	"os"
	"sync"
)

func main() {
	log.Println("starting portal", os.Args)
	ctx, cancel := context.WithCancel(context.Background())
	go osexec.OnShutdown(cancel)

	executable := "portal"
	wait := &sync.WaitGroup{}
	proc := exec.NewRunner[target.App](executable)
	resolve := apps.Resolve(featApps.Path)
	launch := spawn.NewRunner(wait, resolve, proc).Run
	bindings := newRuntimeFactory(ctx, launch)
	run := app.NewRunner(bindings)

	featDispatch := dispatch.NewFeat(executable)
	featServe := serve.NewFeat(launch, tray.NewRunner(launch))
	featOpen := open.NewFeat[target.App](resolve, run)

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

type Adapter struct{ target.Apphost }

func newRuntimeFactory(ctx context.Context, spawn target.Dispatch) target.New {
	invoke := apphost.Invoke(spawn)
	return func(t target.Type, prefix ...string) target.Api {
		switch {
		case t.Is(target.TypeFrontend):
			return &Adapter{Apphost: apphost.NewAdapter(ctx, invoke, prefix...)}
		default:
			return apphost.WithTimeout(ctx, invoke, prefix...)
		}
	}
}
