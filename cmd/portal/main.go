package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/clir"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/dispatch"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	osexec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/resolve"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner"
	"github.com/cryptopunkscc/go-astral-js/runner/app"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
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
	find := resolve.Apps
	spawn := runner.NewSpawner(wait, find, proc)
	bindings := newRuntimeFactory(ctx, spawn)
	run := app.NewRunner(bindings)

	dispatchFeat := dispatch.NewFeat(executable)
	serveFeat := serve.NewFeat(spawn, tray.New(spawn))
	attachFeat := open.NewFeat[target.App](find, run)

	cli := clir.NewCli(ctx, manifest.Name, manifest.Description, version.Run)

	cli.Open(dispatchFeat)
	cli.Serve(serveFeat)
	cli.Attach(attachFeat)
	cli.List(apps.List)
	cli.Install(apps.Install)
	cli.Uninstall(apps.Uninstall)

	err := cli.Run()
	cancel()
	if err != nil {
		log.Println(err)
	}
	wait.Wait()
}

type Adapter struct{ apphost.Flat }

func newRuntimeFactory(ctx context.Context, spawn target.Spawn) target.New {
	invoke := apphost.Invoke(spawn)
	return func(t target.Type, prefix ...string) target.Api {
		switch {
		case t.Is(target.TypeFrontend):
			return &Adapter{Flat: apphost.NewAdapter(ctx, invoke, prefix...)}
		default:
			return apphost.WithTimeout(ctx, invoke, prefix...)
		}
	}
}
