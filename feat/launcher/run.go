package launcher

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
)

func Run(ctx context.Context, bindings runtime.New) (err error) {
	target, err := runner.RawTargetsFS(apps.LauncherSvelteFS)
	if err != nil {
		return
	}
	if len(target) == 0 {
		return errors.New("embed launcher not found")
	}
	return open.RunTarget(ctx, bindings, target[0])
}
