package app

import (
	"context"
	"github.com/cryptopunkscc/portal/target"
)

func Run[T target.Base](run target.Run[T]) target.Run[target.Base] {
	return func(ctx context.Context, app target.Base) error {
		t, ok := app.(T)
		if !ok {
			return target.ErrNotTarget
		}
		return run(ctx, t)
	}
}
