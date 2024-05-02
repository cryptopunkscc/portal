package exec

import (
	"log"
	"time"
)

func Retry[T any](duration time.Duration, fn func(int, int, time.Duration) (T, error)) (t T, err error) {
	t, err = fn(0, 0, 0)
	if err == nil {
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
