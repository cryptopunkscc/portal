package open

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
)

type Deps[T target.Portal_] interface {
	Resolver() target.Resolve[T]
	Runner() target.Run[T]
}

func Runner[T target.Portal_](deps Deps[T]) target.Run[string] {
	resolve := deps.Resolver()
	run := deps.Runner()
	return func(ctx context.Context, path string, args ...string) (err error) {
		if err = apphost.Connect(ctx); err != nil {
			return
		}
		file, err := source.File(path)
		if err != nil {
			return err
		}
		portal, err := resolve(file)
		if err != nil {
			return errors.New("cannot resolve portal: " + err.Error())
		}
		if any(portal) == nil {
			return errors.New("cannot resolve portal for path: " + path)
		}
		plog.Get(ctx).Scope(portal.Manifest().Package).Set(&ctx)
		return run(ctx, portal, args...)
	}
}
