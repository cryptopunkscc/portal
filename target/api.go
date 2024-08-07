package target

import (
	"context"
)

type Api interface{ Apphost }

type NewApi func(context.Context, Portal_) Api

type Path func(src string) (path string, err error)

type File func(path ...string) (source Source, err error)

type Find[T Portal_] func(ctx context.Context, src string) (portals Portals[T], err error)

type Dispatch func(context.Context, string, ...string) (err error)

type Run[T Source] func(ctx context.Context, src T) (err error)

type Runner[T Portal_] interface {
	Run(ctx context.Context, src T) (err error)
	Reload() error
}
