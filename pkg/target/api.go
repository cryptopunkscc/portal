package target

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
)

type Api interface {
	apphost.Flat
}

type New func(p Type, prefix ...string) Api

type Serve func(ctx context.Context, port string, handlers rpc.Handlers, tray Tray) (err error)

type Find[T Portal] func(src string) (portals Portals[T], err error)

type Tray func(ctx context.Context)

type Spawn func(context.Context, string) error

type Run[T Portal] func(ctx context.Context, src T) (err error)

type Lift func(root Source, src string) (result Source, err error)
