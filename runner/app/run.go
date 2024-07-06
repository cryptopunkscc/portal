package app

import (
	"context"
	"github.com/cryptopunkscc/portal/target"
)

func Run[T target.App](run target.Run[T]) target.Run[target.App] {
	return func(ctx context.Context, app target.App) error {
		t, ok := app.(T)
		if !ok {
			return target.ErrNotTarget
		}
		return run(ctx, t)
	}
}
