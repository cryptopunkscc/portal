package launcher

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/prod"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
)

func Run(ctx context.Context, bindings runtime.New) (err error) {
	return prod.RunTargets(ctx, bindings, []runner.Target{
		{Files: apps.LauncherSvelteFS},
		{Files: apps.LauncherBackendFS},
	})
}
