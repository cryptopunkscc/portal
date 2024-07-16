package target

import (
	"context"
	"io/fs"
)

type Api interface{ Apphost }

type NewApi func(context.Context, Portal) Api

type Tray func(ctx context.Context)

type NewTray func(dispatch Dispatch) Tray

type Resolve[T Source] func(src Source) (result T, err error)

type Path func(src string) (path string, err error)

type Finder[T Portal] func(resolve Path, files ...fs.FS) Find[T]

type Find[T Portal] func(ctx context.Context, src string) (portals Portals[T], err error)

type Dispatch func(context.Context, string, ...string) (err error)

type Serve func(ctx context.Context, tray bool) error

type Run[T Portal] func(ctx context.Context, src T) (err error)

type Runner[T Portal] interface {
	Run(ctx context.Context, src T) (err error)
	Reload() error
}

type CreateProject func(Template) error
