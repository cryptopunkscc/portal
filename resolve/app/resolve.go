package app

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/portal"
)

type app[T any] struct{ target.Portal[T] }

func (a *app[T]) IsApp() {}

func Resolve[T any](src target.Source) (t target.App[T], err error) {
	p, err := portal.Resolve[T](src)
	if err == nil {
		t = &app[T]{p}
	}
	return
}
