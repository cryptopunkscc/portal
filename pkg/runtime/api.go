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

type Tray func(ctx context.Context)

type Spawn func(context.Context, string) error

type Run[T target.Portal] func(ctx context.Context, src T) (err error)
