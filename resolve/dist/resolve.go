package dist

import (
	"github.com/cryptopunkscc/portal/resolve/app"
	"github.com/cryptopunkscc/portal/target"
)

type of[T any] struct {
	t T
	target.App[T]
}

func (d *of[T]) Target() T {
	return d.t
}

func (d *of[T]) IsDist() {}

func Resolver[T any](resolve target.Resolve[T]) target.Resolve[target.Dist[T]] {
	return func(src target.Source) (result target.Dist[T], err error) {
		a, err := app.Resolve[T](src)
		if err != nil {
			return
		}
		t, err := resolve(src)
		if err != nil {
			return
		}
		result = &of[T]{App: a, t: t}
		return
	}
}
