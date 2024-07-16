package open

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
)

type Deps[T target.Portal] interface {
	TargetFind() target.Find[T]
	TargetRun() target.Run[T]
}

func Inject[T target.Portal](deps Deps[T]) target.Dispatch {
	return NewFeat(deps.TargetFind(), deps.TargetRun())
}

type Feat[T target.Portal] struct {
	find target.Find[T]
	run  target.Run[T]
}

func NewFeat[T target.Portal](
	find target.Find[T],
	run target.Run[T],
) target.Dispatch {
	return Feat[T]{
		find: find,
		run:  run,
	}.Run
}

func (f Feat[T]) Run(ctx context.Context, path string, _ ...string) (err error) {
	plog.Get(ctx).Type(f).Set(&ctx)
	portal, err := f.find(ctx, path)
	if err != nil {
		return errors.New("cannot resolve portal: " + err.Error())
	}
	for _, t := range portal {
		plog.Get(ctx).Scope(t.Manifest().Package).Set(&ctx)
		return f.run(ctx, t)
	}
	return errors.New("no target found")
}
