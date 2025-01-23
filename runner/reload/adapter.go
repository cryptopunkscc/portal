package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runtime/bind"
)

func Adapter[T target.Portal_](runner target.ReRunner[T]) func(bind.NewRuntime) target.ReRunner[T] {
	return func(newRuntime bind.NewRuntime) target.ReRunner[T] {
		return adapter[T]{
			newRuntime: newRuntime,
			inner:      runner,
		}
	}
}

type adapter[T target.Portal_] struct {
	newRuntime bind.NewRuntime
	inner      target.ReRunner[T]
}

func (a adapter[T]) Run(ctx context.Context, src T, args ...string) (err error) {
	a.newRuntime(ctx, src)
	return a.inner.Run(ctx, src, args...)
}

func (a adapter[T]) ReRun() error {
	return a.inner.ReRun()
}
