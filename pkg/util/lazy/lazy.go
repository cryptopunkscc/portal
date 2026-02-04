package lazy

import "sync"

func V[T any](create func() T) func() T {
	var instance *T
	return func() T {
		if instance == nil {
			i := create()
			instance = &i
		}
		return *instance
	}
}

func S[T any](new func() T) func() T {
	var mu sync.Mutex
	return V(func() T {
		mu.Lock()
		defer mu.Unlock()
		return new()
	})
}
