package apps

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/mock/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
)

func Observe(ctx context.Context, conn rpc.Conn) (err error) {
	plog.Get(ctx).Scope("apps.Observe").Set(&ctx).Println(conn)
	return appstore.Observe(ctx, conn)
}
