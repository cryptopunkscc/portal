package target

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Run[T any] func(ctx context.Context, src T, args ...string) (err error)

func (r Run[T]) Run(ctx context.Context, src T, args ...string) (err error) {
	return r(ctx, src, args...)
}

func (r Run[T]) Start(ctx context.Context, src T, args ...string) (err error) {
	go func() {
		if err = r(ctx, src, args...); err != nil {
			plog.Get(ctx).Type(r).Println("Start failed:", err)
		}
	}()
	return nil
}

func (r Run[T]) Portal() Run[Portal_] {
	return func(ctx context.Context, src Portal_, args ...string) error {
		if t, ok := src.(T); ok {
			return r(ctx, t, args...)
		}
		return ErrNotTarget
	}
}
