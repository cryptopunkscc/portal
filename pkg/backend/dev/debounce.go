package backend

import "time"

func debounce[T any](in <-chan T, t time.Duration) (out <-chan T) {
	o := make(chan T)
	out = o
	go func() {
		defer close(o)
		last := int64(0)
		threshold := t.Nanoseconds()
		for l := range in {
			current := time.Now().UnixNano()
			if current-last > threshold {
				o <- l
				last = current
			}
		}
	}()
	return
}
