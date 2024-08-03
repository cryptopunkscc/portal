package target

import (
	"context"
	"fmt"
)

func (r Run[T]) Portal() Run[Base] {
	return func(ctx context.Context, src Base) (err error) {
		t, ok := src.(T)
		if !ok {
			return fmt.Errorf("cannot run %T as %T", src, t)
		}
		return r(ctx, t)
	}
}
