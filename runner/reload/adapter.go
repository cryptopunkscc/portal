package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
)

func Adapter[T target.Portal_](runner target.ReRunner[T]) func(bind.NewCore) target.ReRunner[T] {
	return func(newCore bind.NewCore) target.ReRunner[T] {
		return adapter[T]{
			newCore: newCore,
			inner:   runner,
		}
	}
}

type adapter[T target.Portal_] struct {
	newCore bind.NewCore
	inner   target.ReRunner[T]
}

func (a adapter[T]) Run(ctx context.Context, src T, args ...string) (err error) {
	a.newCore(ctx, src)
	return a.inner.Run(ctx, src, args...)
}

func (a adapter[T]) ReRun() error {
	return a.inner.ReRun()
}
