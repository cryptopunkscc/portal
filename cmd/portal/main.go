package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/clir"
	"github.com/cryptopunkscc/go-astral-js/deps/dispatch"
	"github.com/cryptopunkscc/go-astral-js/deps/find"
	"github.com/cryptopunkscc/go-astral-js/deps/open"
	"github.com/cryptopunkscc/go-astral-js/deps/serve"
	featApps "github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	osexec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"os"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osexec.OnShutdown(cancel)

	println("...")
	log := plog.New().I().Set(&ctx)
	log.Scope("main").Println("starting portal", os.Args)

	executable := "portal"
	port := "portal"
	portOpen := "portal.open"
	wait := &sync.WaitGroup{}
	targetCache := target.NewCache[target.App]()
	findApps := find.Create(targetCache, apps.NewFind)
	featDispatch := dispatch.Create(executable, portOpen)
	featServe := serve.Create(wait, executable, port, findApps)
	featOpen := open.Create(portOpen, NewAdapter, findApps)
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
