package bundle

import (
	"github.com/cryptopunkscc/portal/target2"
)

type of[T any] struct{ target2.Dist[T] }

func (t of[T]) IsBundle() {}

func Resolver[T any](resolve target2.Resolve[target2.Dist[T]]) target2.Resolve[target2.AppBundle[T]] {
	return func(src target2.Source) (app target2.AppBundle[T], err error) {
		b, err := Resolve(src)
		if err != nil {
			return
		}
		td := &of[T]{}
		if td.Dist, err = resolve(b); err != nil {
			return
		}
		app = td
		return
	}
}
