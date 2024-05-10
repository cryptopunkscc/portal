package runtime

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

type Api interface {
	apphost.Flat
}

type New func(p target.Type, prefix ...string) Api

type Serve func(ctx context.Context, port string, handlers rpc.Handlers, tray Tray) (err error)

type Resolve[T target.Portal] func(src string) (portals target.Portals[T], err error)

type Attach[T target.Portal] func(ctx context.Context, bindings New, app T, prefix ...string) (err error)

type Tray func(ctx context.Context)
