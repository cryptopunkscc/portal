package launcher

import (
	"github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/prod"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
)

func Run(bindings runner.Bindings) (err error) {
	return prod.RunTargets(bindings, []runner.Target{
		{Files: apps.LauncherSvelteFS},
		{Files: apps.LauncherBackendFS},
	})
}
