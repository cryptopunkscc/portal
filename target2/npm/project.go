package npm

import (
	"github.com/cryptopunkscc/portal/target2"
)

type project[T any] struct {
	target2.NodeModule
	target2.Portal[T]
	resolveDist target2.Resolve[target2.Dist[T]]
}

func (p project[T]) IsProject() {}
func (p project[T]) Dist() (t target2.Dist[T]) {
	sub, err := p.Sub("dist")
	if err != nil {
		return
	}
	t, err = p.resolveDist(sub)
	return
}

func Resolver[T any](resolve target2.Resolve[target2.Dist[T]]) target2.Resolve[target2.Project[T]] {
	return func(src target2.Source) (t target2.Project[T], err error) {
		p := &project[T]{}
		if p.NodeModule, err = Resolve(src); err != nil {
			return
		}
		if p.Portal, err = resolve(src); err != nil {
			return
		}
		t = p
		return
	}
}
