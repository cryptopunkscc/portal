package project

import (
	json2 "encoding/json"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/cryptopunkscc/portal/resolve/dist"
)

func Resolver[T any](resolve target.Resolve[T]) target.Resolve[target.Project[T]] {
	return func(src target.Source) (t target.Project[T], err error) {
		if _, err = resolve(src); err != nil {
			return
		}
		p := &project[T]{}
		if err = all.Unmarshalers.Load(&p.manifest, src.FS(), target.BuildFilename); err != nil {
			return
		}
		p.build = target.LoadBuilds(src)
		p.resolveDist = dist.Resolver(resolve)
		p.Source = src
		if p.manifest.Exec == "" {
			p.manifest.Exec = target.GetBuild(p).Exec
		}
		t = p
		return
	}
}

type project[T any] struct {
	target.Source
	build       target.Builds
	manifest    target.Manifest
	resolveDist target.Resolve[target.Dist[T]]
}

func (p *project[T]) Changed(skip ...string) bool  { return target.Changed(p, skip...) }
func (p *project[T]) MarshalJSON() ([]byte, error) { return json2.Marshal(p.Manifest()) }
func (p *project[T]) Manifest() *target.Manifest   { return &p.manifest }
func (p *project[T]) Target() T                    { return p.Dist().Target() }
func (p *project[T]) Build() target.Builds         { return p.build }
func (p *project[T]) Dist_() (t target.Dist_)      { return p.Dist() }
func (p *project[T]) Dist() (t target.Dist[T]) {
	sub, err := p.Sub("dist")
	if err != nil {
		panic(err)
	}
	t, err = p.resolveDist(sub)
	if err != nil {
		panic(err)
	}
	return
}

var Resolve_ = Resolver(target.Resolve_)
