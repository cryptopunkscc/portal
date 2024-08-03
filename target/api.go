package target

import (
	"context"
)

type Api interface{ Apphost }

type NewApi func(context.Context, Base) Api

type Tray func(ctx context.Context) error

type NewTray func(dispatch Dispatch) Tray

type Resolve[T any] func(src Source) (result T, err error)

type Path func(src string) (path string, err error)

type Find[T Base] func(ctx context.Context, src string) (portals Portals[T], err error)

type Dispatch func(context.Context, string, ...string) (err error)

type Serve func(ctx context.Context, tray bool) error

type Run[T Base] func(ctx context.Context, src T) (err error)

type Runner[T Base] interface {
	Run(ctx context.Context, src T) (err error)
	Reload() error
}

type CreateProject func(Template) error
