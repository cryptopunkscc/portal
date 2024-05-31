package flow

import "time"

type Input[T any] <-chan T

func From[T any](in <-chan T) Input[T] {
	return in
}

func (f Input[T]) Debounce(t time.Duration) (input Input[T]) {
	o := make(chan T)
	input = o
	go func() {
		defer close(o)
		last := int64(0)
		threshold := t.Nanoseconds()
		for l := range f {
			current := time.Now().UnixNano()
			if current-last > threshold {
				o <- l
				last = current
			}
		}
	}()
	return
}

func Map[T1 any, T2 any](from <-chan T1, transform func(T1) (T2, bool)) (to Input[T2]) {
	c := make(chan T2)
	to = c
	go func() {
		for v1 := range from {
			if v2, ok := transform(v1); ok {
				c <- v2
			}
		}
	}()
	return
}
