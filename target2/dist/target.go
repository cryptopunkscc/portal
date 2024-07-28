package dist

import (
	"github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/app"
)

type of[T any] struct {
	t T
	target2.App[T]
}

func (d *of[T]) Target() T {
	return d.t
}

func (d *of[T]) IsDist() {}

func Resolver[T any](resolve target2.Resolve[T]) target2.Resolve[target2.Dist[T]] {
	return func(src target2.Source) (result target2.Dist[T], err error) {
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
