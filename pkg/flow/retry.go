package flow

import (
	"context"
	"time"
)

func Retry(
	ctx context.Context,
	duration time.Duration,
	retry func(i int, n int, d time.Duration) (err error),
) (err error) {
	retryT := func(i int, n int, d time.Duration) (_ any, err error) {
		return nil, retry(i, n, d)
	}
	_, err = RetryT[any](ctx, duration, retryT)
	return
}

func RetryT[T any](
	ctx context.Context,
	duration time.Duration,
	retry func(i int, n int, d time.Duration) (t T, err error),
) (t T, err error) {
	if t, err = retry(0, 0, 0); err == nil {
		return
	}
	retries := AwaitExp(duration)
	n := len(retries)
	for i, d := range retries {
		time.Sleep(d)
		if ctx.Err() != nil {
			return
		}
		t, err = retry(i+1, n, d)
		if err == nil {
			return
		}
	}
	return
}
