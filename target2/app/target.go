package app

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/portal"
)

type app[T any] struct{ target.Portal[T] }

func (a *app[T]) IsApp() {}

func New[T any](src target.Portal[T]) (t target.App[T], err error) {
	t = &app[T]{src}
	return
}

func Resolve[T any](src target.Source) (t target.App[T], err error) {
	p, err := portal.Resolve[T](src)
	if err == nil {
		t = &app[T]{p}
	}
	return
}
