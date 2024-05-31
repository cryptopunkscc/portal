package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	manifest "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/clir"
	feature "github.com/cryptopunkscc/go-astral-js/feat"
	featApps "github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/feat/create"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
	"github.com/cryptopunkscc/go-astral-js/pkg/broadcast"
	osExec "github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	create2 "github.com/cryptopunkscc/go-astral-js/runner/create"
	"github.com/cryptopunkscc/go-astral-js/runner/dev"
	"github.com/cryptopunkscc/go-astral-js/runner/dist"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	"github.com/cryptopunkscc/go-astral-js/runner/pack"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/runner/service"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/portals"
	"os"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go osExec.OnShutdown(cancel)

	println("...")
	plog.ErrorStackTrace = true
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")

	scope := feature.Scope[target.Portal]{
		Port:           "dev.portal",
		Prefix:         []string{"dev"},
		WrapApi:        NewAdapter,
		WaitGroup:      &sync.WaitGroup{},
		TargetCache:    target.NewCache[target.Portal](),
		NewRunTarget:   dev.NewRun,
		NewRunService:  service.NewRun,
		TargetFinder:   portals.NewFind,
		ExecTarget:     exec.NewRun[target.Portal]("portal-dev"),
		GetPath:        featApps.Path,
		FeatObserve:    featApps.Observe,
		JoinTarget:     query.NewRunner[target.App]("dev.portal.open").Run,
		DispatchTarget: query.NewRunner[target.App]("dev.portal.open").Start,
		Processes:      &sig.Map[string, target.Portal]{},
	}
	scope.RpcHandlers = rpc.Handlers{
		"ctrl": broadcast.New(scope.GetProcesses(), "ctrl", scope.Prefix...).Signal,
	}
	scope.DispatchService = scope.GetServeFeature().Dispatch

	featBuild := build.NewFeat(dist.NewRun, pack.Run)
	featCreate := create.NewFeat(create2.NewRun, featBuild.Dist).Run

	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(scope.GetDispatchFeature())
	cli.Open(scope.GetOpenFeature())
	cli.Create(create.List, featCreate)
	cli.Build(featBuild.Run)
	cli.Portals(scope.GetTargetFind())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	scope.WaitGroup.Wait()
}

type Adapter struct{ target.Api }

func NewAdapter(api target.Api) target.Api { return &Adapter{Api: api} }
