package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/target"
)

func Adapter[T target.Portal_](runner target.Runner[T]) func(bind.NewRuntime) target.Runner[T] {
	return func(newRuntime bind.NewRuntime) target.Runner[T] {
		return adapter[T]{
			newRuntime: newRuntime,
			inner:      runner,
		}
	}
}

type adapter[T target.Portal_] struct {
	newRuntime bind.NewRuntime
	inner      target.Runner[T]
}

func (a adapter[T]) Run(ctx context.Context, src T) (err error) {
	a.newRuntime(ctx, src)
	return a.inner.Run(ctx, src)
}

func (a adapter[T]) Reload() error {
	return a.inner.Reload()
}
