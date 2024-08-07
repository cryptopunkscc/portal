package target

import (
	"context"
)

func (r Run[T]) Call(ctx context.Context, src T) (err error) {
	return r(ctx, src)
}

func (r Run[T]) Portal(ctx context.Context, src Portal_) (err error) {
	t, ok := src.(T)
	if !ok {
		return ErrNotTarget
	}
	return r(ctx, t)
}
