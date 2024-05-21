package target

import (
	"context"
)

type Api interface {
	Apphost
}

type NewApi func(context.Context, Portal) Api

type Tray func(ctx context.Context)

type Path func(src string) (path string, err error)

type Dispatch func(context.Context, string) (err error)

type Resolve[T Source] func(src Source) (result T, err error)

type Find[T Portal] func(src string) (portals Portals[T], err error)

type Finder[T Portal] func(resolve Path) Find[T]

type Run[T Portal] func(ctx context.Context, src T) (err error)
