package target

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Run[T any] func(ctx context.Context, src T, args ...string) (err error)

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

func (r Run[T]) Start(ctx context.Context, src T, args ...string) (err error) {
	go func() {
		if err = r(ctx, src, args...); err != nil {
			plog.Get(ctx).Type(r).Println("Start failed:", err)
		}
	}()
	return nil
}
