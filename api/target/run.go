package target

import (
	"context"
)

type Run[T any] func(ctx context.Context, src T, args ...string) (err error)

type Runner[T Portal_] interface {
	Run(ctx context.Context, src T, args ...string) (err error)
	Reload() error
}

func (r Run[T]) Call(ctx context.Context, src T, args ...string) (err error) {
	return r(ctx, src, args...)
}

func (r Run[T]) Portal(ctx context.Context, src Portal_, args ...string) (err error) {
	t, ok := src.(T)
	if !ok {
		return ErrNotTarget
	}
	return r(ctx, t, args...)
}
