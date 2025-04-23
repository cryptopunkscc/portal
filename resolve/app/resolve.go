package app

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/portal"
)

type of[T any] struct{ target.Portal[T] }

func (a *of[T]) IsApp() {}

func Resolve[T any](src target.Source) (t target.App[T], err error) {
	p, err := portal.Resolve[T](src)
	if err == nil {
		t = &of[T]{p}
	}
	return
}
