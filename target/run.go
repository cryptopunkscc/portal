package target

import (
	"context"
)

type Run[T Source] func(ctx context.Context, src T) (err error)

type Runner[T Portal_] interface {
	Run(ctx context.Context, src T) (err error)
	Reload() error
}

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
