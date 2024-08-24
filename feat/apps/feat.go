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
	return rpc.
		NewApp("portal.apps").
		Client(deps.Client()).
		Interface(deps.Apps()).
		Logger(plog.New()).
		Routes("*").
		Run
}
