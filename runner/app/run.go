package app

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
)

func Run[T target.Portal_](run target.Run[T]) target.Run[target.Portal_] {
	return func(ctx context.Context, app target.Portal_) error {
		t, ok := app.(T)
		if !ok {
			return target.ErrNotTarget
		}
		return run(ctx, t)
	}
}
