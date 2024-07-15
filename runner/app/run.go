package app

import (
	"context"
	"github.com/cryptopunkscc/portal/target"
)

func Run[T target.Portal](run target.Run[T]) target.Run[target.Portal] {
	return func(ctx context.Context, app target.Portal) error {
		t, ok := app.(T)
		if !ok {
			return target.ErrNotTarget
		}
		return run(ctx, t)
	}
}
