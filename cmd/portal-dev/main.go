package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/clir"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/feat/create"
	dev "github.com/cryptopunkscc/go-astral-js/feat/dev"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/templates"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	osexec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/resolve"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner"
	dev2 "github.com/cryptopunkscc/go-astral-js/runner/dev"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"log"
	"os"
	"sync"
)

func main() {
	log.Println("starting portal development", os.Args)
	ctx, cancel := context.WithCancel(context.Background())
	go osexec.OnShutdown(cancel)

	wait := &sync.WaitGroup{}
	proc := exec.NewRunner[target.Portal]("portal-dev")
	find := resolve.Portals(appstore.Path)
	spawn := runner.NewSpawner(wait, find, proc)
	bindings := newRuntimeFactory(ctx, spawn)
	run := dev2.NewRunner(bindings)

	devFeat := dev.NewFeat(wait, spawn)
	attachFeat := open.NewFeat[target.Portal](find, run)

	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)

	cli.Dev(devFeat)
	cli.Attach(attachFeat)
	cli.Create(templates.List, create.Run)
	cli.Build(build.Run)
	cli.Apps()

	err := cli.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("closing portal development")
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
