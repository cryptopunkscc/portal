package exec

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"time"
)

func Retry(ctx context.Context, duration time.Duration, fn func(int, int, time.Duration) error) (err error) {
	_, err = RetryT[any](ctx, duration, func(i int, i2 int, duration time.Duration) (_ any, err error) {
		err = fn(i2, i, duration)
		return
	})
	return
}

func RetryT[T any](ctx context.Context, duration time.Duration, fn func(int, int, time.Duration) (T, error)) (t T, err error) {
	log := plog.Get(ctx)
	if t, err = fn(0, 0, 0); err == nil {
		return
	}
	retries := AwaitExp(duration)
	n := len(retries)
	for i, d := range retries {
		log.Printf("%d/%d attempt %v: retry after %v", i+1, n, err, d)
		time.Sleep(d)
		if ctx.Err() != nil {
			return
		}
		t, err = fn(i+1, n, d)
		if err == nil {
			return
		}
	}
	return
}
