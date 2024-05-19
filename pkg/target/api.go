package target

import (
	"context"
)

type Api interface {
	Apphost
}

type New func(p Type, prefix ...string) Api

type Tray func(ctx context.Context)

type Path func(src string) (string, error)

type Dispatch func(context.Context, string) error

type ResolveT[T Source] func(src Source) (result T, err error)

type Find[T Portal] func(src string) (portals Portals[T], err error)

type Run[T Portal] func(ctx context.Context, src T) (err error)
