package dist

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/app"
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
		t, err := resolve(a)
		if err != nil {
			return
		}
		result = &of[T]{App: a, t: t}
		return
	}
}

var ResolveAny = Resolver[any](func(target.Source) (result any, err error) { return })
