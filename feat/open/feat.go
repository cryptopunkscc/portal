package open

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
)

type Deps[T target.Portal_] interface {
	TargetResolve() target.Resolve[T]
	TargetRun() target.Run[T]
}

func Inject[T target.Portal_](deps Deps[T]) target.Dispatch {
	return NewFeat(deps.TargetResolve(), deps.TargetRun())
}

type Feat[T target.Portal_] struct {
	resolve target.Resolve[T]
	run     target.Run[T]
}

func NewFeat[T target.Portal_](
	resolve target.Resolve[T],
	run target.Run[T],
) target.Dispatch {
	return Feat[T]{
		resolve: resolve,
		run:     run,
	}.Run
}

func (f Feat[T]) Run(ctx context.Context, path string, _ ...string) (err error) {
	plog.Get(ctx).Type(f).Set(&ctx)
	file, err := source.File(path)
	if err != nil {
		return err
	}
	portal, err := f.resolve(file)
	if err != nil {
		return errors.New("cannot resolve portal: " + err.Error())
	}
	if any(portal) == nil {
		return errors.New("cannot resolve portal for path: " + path)
	}
	plog.Get(ctx).Scope(portal.Manifest().Package).Set(&ctx)
	return f.run(ctx, portal)
}
