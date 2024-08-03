package bundle

import "github.com/cryptopunkscc/portal/target"

type of[T any] struct{ target.Dist[T] }

func (t of[T]) IsBundle() {}

func Resolver[T any](resolve target.Resolve[target.Dist[T]]) target.Resolve[target.AppBundle[T]] {
	return func(src target.Source) (app target.AppBundle[T], err error) {
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
