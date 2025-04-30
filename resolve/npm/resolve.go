package npm

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/project"
)

var Resolve_ = Resolver(target.Resolve_)

func Resolver[T any](resolve target.Resolve[T]) target.Resolve[target.ProjectNpm[T]] {
	r := project.Resolver(resolve)
	return func(src target.Source) (t target.ProjectNpm[T], err error) {
		p := &Project[T]{}
		if p.Project, err = r.Resolve(src); err != nil {
			return
		}
		if p.nodeModule, err = ResolveNodeModule(src); err != nil {
			return
		}
		t = p
		return
	}
}

func ResolveNodeModule(src target.Source) (t target.NodeModule, err error) {
	if !src.IsDir() {
		return nil, target.ErrNotTarget
	}
	s := &Module{Source: src}
	if err = s.LoadPkgJson(); err != nil {
		return
	}
	t = s
	return
}
