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
	"github.com/cryptopunkscc/go-astral-js/mock/appstore"
	osExec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	devRunner "github.com/cryptopunkscc/go-astral-js/runner/dev"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/runner/spawn"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apphost"
	"github.com/cryptopunkscc/go-astral-js/target/portals"
	"log"
	"os"
	"sync"
)

func main() {
	log.Println("starting portal development", os.Args)
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	wait := &sync.WaitGroup{}
	executable := "portal-dev"
	prefix := "dev"

	findPortals := target.Cached(portals.Find)(appstore.Path)

	runQuery := query.NewRunner[target.Portal](prefix).Run
	newApphost := apphost.NewFactory(runQuery, prefix)
	newApi := target.ApiFactory(
		NewAdapter,
		newApphost.NewAdapter,
		newApphost.WithTimeout,
	)
	runDev := devRunner.NewRunner(newApi)
	runProc := exec.NewRunner[target.Portal](executable)
	runSpawn := spawn.NewRunner(wait, findPortals, runProc).Run

	featDev := dev.NewFeat("dev.portal", wait, runSpawn, runQuery)
	featOpen := open.NewFeat[target.Portal](findPortals, runDev)
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
	wait.Wait()
	log.Println("closing portal development")
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }
