package app

import (
	"github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/portal"
)

type app[T any] struct{ target2.Portal[T] }

func (a *app[T]) IsApp() {}

func New[T any](src target2.Portal[T]) (t target2.App[T], err error) {
	t = &app[T]{src}
	return
}

func Resolve[T any](src target2.Source) (t target2.App[T], err error) {
	p, err := portal.Resolve[T](src)
	if err == nil {
		t = &app[T]{p}
	}
	return
}
