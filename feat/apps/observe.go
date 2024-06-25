package apps

import (
	"context"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Observe(ctx context.Context, conn rpc.Conn) (err error) {
	plog.Get(ctx).Scope("apps.Observe").Set(&ctx).Println(conn)
	return appstore.Observe(ctx, conn)
}
