package main

import (
	"context"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/builder"
	"github.com/cryptopunkscc/go-astral-js/clir"
	featApps "github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/feat/create"
	serve2 "github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	osExec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	devRunner "github.com/cryptopunkscc/go-astral-js/runner/dev"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/portals"
	"os"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	println("...")
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")

	scope := builder.Scope[target.Portal]{
		Port:           "dev.portal",
		Prefix:         []string{"dev"},
		WrapApi:        NewAdapter,
		WaitGroup:      &sync.WaitGroup{},
		TargetCache:    target.NewCache[target.Portal](),
		NewTargetRun:   devRunner.NewRun,
		NewServe:       serve.NewRun,
		TargetFinder:   portals.NewFind,
		ExecTarget:     exec.NewRun[target.Portal]("portal-dev"),
		AppsPath:       featApps.Path,
		FeatObserve:    featApps.Observe,
		DispatchTarget: query.NewRunner[target.App]("dev.portal.open").Run,
	}

	scope.DispatchService = func(ctx context.Context, _ string, _ ...string) (err error) {
		srv := scope.GetServeFeature()
		go func() {
			if err = srv(ctx, false); err != nil {
				plog.Get(ctx).Type(serve2.Feat{}).Println(err)
			}
		}()
		return
	}

	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(scope.GetDispatchFeature())
	cli.Open(scope.GetOpenFeature())
	cli.Create(create.List, create.NewFeat().Run)
	cli.Build(build.NewFeat().Run)
	cli.Portals(scope.GetTargetFind())

	err := cli.Run()
	if err != nil {
		log.Println(err)
	}
	scope.WaitGroup.Wait()
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }
