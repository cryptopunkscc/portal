package launcher

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
)

func Run(ctx context.Context, bindings runtime.New) (err error) {
	if noService := rpc.Command(portal.Request, "ping") != nil; noService {
		if err = portal.Open(nil, "serve", "-t").Start(); err != nil {
			return
		}
	}
	target, ok := <-project.RawTargets(apps.LauncherSvelteFS)
	if !ok {
		return errors.New("embed launcher not found")
	}
	return open.RunTarget(ctx, bindings, target)
}
