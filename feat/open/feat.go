package open

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
)

type Deps[T target.Portal_] interface {
	Resolver() target.Resolve[T]
	Runner() target.Run[T]
}

func Feat[T target.Portal_](
	deps Deps[T],
) target.Request {
	resolve := deps.Resolver()
	run := deps.Runner()
	return func(ctx context.Context, path string) (err error) {
		plog.Get(ctx)
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
		return run(ctx, portal)
	}
}
