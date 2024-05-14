package apps

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
)

func Observe(ctx context.Context, conn rpc.Conn) (err error) {
	return appstore.Observe(ctx, conn)
}
