package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/clir"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/feat/create"
	"github.com/cryptopunkscc/go-astral-js/feat/dev"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	osExec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/portals"
	devRunner "github.com/cryptopunkscc/go-astral-js/runner/dev"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/spawn"
	"log"
	"os"
	"sync"
)

func main() {
	log.Println("starting portal development", os.Args)
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	wait := &sync.WaitGroup{}
	proc := exec.NewRunner[target.Portal]("portal-dev")
	find := portals.Resolve(appstore.Path)
	launch := spawn.NewRunner(wait, find, proc).Run
	bindings := newRuntimeFactory(ctx, launch)
	run := devRunner.NewRunner(bindings)

	featDev := dev.NewFeat(wait, launch)
	featOpen := open.NewFeat[target.Portal](find, run)
	featBuild := build.NewFeat().Run
	featCreate := create.NewFeat().Run

	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)

	cli.Dev(featDev)
	cli.Open(featOpen)
	cli.Create(create.List, featCreate)
	cli.Build(featBuild)
	cli.Apps()

	err := cli.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("closing portal development")
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
