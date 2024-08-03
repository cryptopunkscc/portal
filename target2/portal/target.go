package portal

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/base"
)

type portal[T any] struct {
	target.Base
}

func (a *portal[T]) IsApp()        {}
func (a *portal[T]) Target() (t T) { return }

func New[T any](src target.Base) (t target.Portal[T], err error) {
	t = &portal[T]{src}
	return
}

func Resolve[T any](src target.Source) (t target.Portal[T], err error) {
	b, err := base.ResolveBase(src)
	if err != nil {
		return
	}
	t = &portal[T]{b}
	return
}
