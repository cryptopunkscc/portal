package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/target"
)

func Adapter[T target.Portal_](runner target.Runner[T]) func(target.NewApi) target.Runner[T] {
	return func(newApi target.NewApi) target.Runner[T] {
		return adapter[T]{
			newApi: newApi,
			inner:  runner,
		}
	}
}

type adapter[T target.Portal_] struct {
	newApi target.NewApi
	inner  target.Runner[T]
}

func (a adapter[T]) Run(ctx context.Context, src T) (err error) {
	a.newApi(ctx, src)
	return a.inner.Run(ctx, src)
}

func (a adapter[T]) Reload() error {
	return a.inner.Reload()
}
