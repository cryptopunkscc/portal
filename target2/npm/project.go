package npm

import (
	"github.com/cryptopunkscc/portal/target"
)

type project[T any] struct {
	nodeModule target.NodeModule
	target.Portal[T]
	resolveDist target.Resolve[target.Dist[T]]
}

func (p project[T]) IsProject()                   {}
func (p project[T]) PkgJson() *target.PackageJson { return p.nodeModule.PkgJson() }
func (p project[T]) Dist_() (t target.Dist_)      { return p.Dist() }
func (p project[T]) Dist() (t target.Dist[T]) {
	sub, err := p.Sub("dist")
	if err != nil {
		return
	}
	t, err = p.resolveDist(sub)
	return
}

func Resolver[T any](resolve target.Resolve[target.Dist[T]]) target.Resolve[target.ProjectNpm[T]] {
	return func(src target.Source) (t target.ProjectNpm[T], err error) {
		p := &project[T]{}
		if p.nodeModule, err = Resolve(src); err != nil {
			return
		}
		if p.Portal, err = resolve(src); err != nil {
			return
		}
		p.resolveDist = resolve
		t = p
		return
	}
}
