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
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	osexec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
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
	find := portal.Resolve(appstore.Path)
	launch := spawn.NewRunner(wait, find, proc)
	bindings := newRuntimeFactory(ctx, launch)
	run := app.NewRunner(bindings)

	dispatchFeat := dispatch.NewFeat(executable)
	serveFeat := serve.NewFeat(launch, tray.NewRunner(launch))
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

type Adapter struct{ target.Apphost }

func newRuntimeFactory(ctx context.Context, spawn target.Spawn) target.New {
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
