package portal

import (
	"github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/base"
)

type portal[T any] struct {
	target2.Base
}

func (a *portal[T]) IsApp()        {}
func (a *portal[T]) Target() (t T) { return }

func New[T any](src target2.Base) (t target2.Portal[T], err error) {
	t = &portal[T]{src}
	return
}

func Resolve[T any](src target2.Source) (t target2.Portal[T], err error) {
	b, err := base.ResolveBase(src)
	if err != nil {
		return
	}
	t = &portal[T]{b}
	return
}
