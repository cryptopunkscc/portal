package apps

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/apps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

type Serve func(ctx context.Context) error

type Deps interface {
	Apps() apps.Apps
	Client() apphost.Client
}

func Feat(deps Deps) Serve {
	router := rpc.
		NewApp("cc.cryptopunks.portal.apps").
		Client(deps.Client()).
		Interface(deps.Apps()).
		Routes("*")
	return func(ctx context.Context) error {
		return router.
			Logger(plog.Get(ctx).Scope("apps")).
			Run(ctx)
	}
}
