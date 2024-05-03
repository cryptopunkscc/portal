package exec

import (
	"log"
	"time"
)

func Retry(duration time.Duration, fn func(int, int, time.Duration) error) (err error) {
	_, err = RetryT[any](duration, func(i int, i2 int, duration time.Duration) (_ any, err error) {
		err = fn(i2, i, duration)
		return
	})
	return
}

func RetryT[T any](duration time.Duration, fn func(int, int, time.Duration) (T, error)) (t T, err error) {
	if t, err = fn(0, 0, 0); err == nil {
		return
	}
	retries := AwaitExp(duration)
	n := len(retries)
	for i, d := range retries {
		log.Printf("%d/%d attempt %v: retry after %v\n", i+1, n, err, d)
		time.Sleep(d)
		t, err = fn(i+1, n, d)
		if err == nil {
			return
		}
	}
	return
}
