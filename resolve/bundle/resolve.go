package bundle

import (
	"github.com/cryptopunkscc/portal/resolve/zip"
	"github.com/cryptopunkscc/portal/target"
)

type of[T any] struct {
	target.App[T]
	bundle target.Bundle
}

func (t of[T]) Package() target.Source { return t.bundle.Package() }

func Resolver[T any](resolve target.Resolve[target.Dist[T]]) target.Resolve[target.AppBundle[T]] {
	return func(src target.Source) (app target.AppBundle[T], err error) {
		b, err := zip.Resolve(src)
		if err != nil {
			return
		}
		td := &of[T]{}
		if td.App, err = resolve(b); err != nil {
			return
		}
		td.bundle = b
		app = td
		return
	}
}
