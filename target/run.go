package target

import (
	"context"
)

func (r Run[T]) Portal() Run[Portal_] {
	return func(ctx context.Context, src Portal_) (err error) {
		t, ok := src.(T)
		if !ok {
			return ErrNotTarget
		}
		return r(ctx, t)
	}
}
