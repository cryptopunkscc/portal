package serve

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/runner/exec"
	runnerServe "github.com/cryptopunkscc/go-astral-js/runner/serve"
	"github.com/cryptopunkscc/go-astral-js/runner/spawn"
	"github.com/cryptopunkscc/go-astral-js/runner/tray"
	"github.com/cryptopunkscc/go-astral-js/target"
	"sync"
)

func Create(
	wait *sync.WaitGroup,
	executable string,
	port string,
	findApps target.Find[target.App],
) func(context.Context, bool) error {
	runProc := exec.NewRun[target.App](executable)
	runSpawn := spawn.NewRunner(wait, findApps, runProc).Run
	return serve.NewFeat(
		port,
		runSpawn,
		tray.New(runSpawn),
		runnerServe.NewRun,
		apps.Observe,
		apps.Install,
		apps.Uninstall,
	)
}
