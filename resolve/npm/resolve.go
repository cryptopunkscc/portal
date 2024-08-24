package npm

import (
	json2 "encoding/json"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/cryptopunkscc/portal/pkg/dec/json"
	"github.com/cryptopunkscc/portal/resolve/dist"
)

type nodeModule struct {
	target.Source
	packageJson *target.PackageJson
}

func (n *nodeModule) PkgJson() *target.PackageJson {
	return n.packageJson
}

func (n *nodeModule) LoadPkgJson() error {
	return json.Unmarshaler.Load(&n.packageJson, n.Files(), target.PackageJsonFilename)
}

func Resolve(src target.Source) (t target.NodeModule, err error) {
	if !src.IsDir() {
		return nil, target.ErrNotTarget
	}
	s := &nodeModule{Source: src}
	if err = s.LoadPkgJson(); err != nil {
		return
	}
	t = s
	return
}

type project[T any] struct {
	target.Source
	build       target.Builds
	nodeModule  target.NodeModule
	manifest    target.Manifest
	resolveDist target.Resolve[target.Dist[T]]
}

func (p *project[T]) MarshalJSON() ([]byte, error) { return json2.Marshal(p.Manifest()) }
func (p *project[T]) Manifest() *target.Manifest   { return &p.manifest }
func (p *project[T]) Target() T                    { return p.Dist().Target() }
func (p *project[T]) Build() target.Builds         { return p.build }
func (p *project[T]) PkgJson() *target.PackageJson { return p.nodeModule.PkgJson() }
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

func Resolver[T any](resolve target.Resolve[T]) target.Resolve[target.ProjectNpm[T]] {
	return func(src target.Source) (t target.ProjectNpm[T], err error) {
		if _, err = resolve(src); err != nil {
			return
		}
		p := &project[T]{}
		if err = all.Unmarshalers.Load(&p.manifest, src.Files(), target.BuildFilename); err != nil {
			return
		}
		if p.nodeModule, err = Resolve(src); err != nil {
			return
		}
		p.build = target.LoadBuilds(src)
		p.resolveDist = dist.Resolver[T](resolve)
		p.Source = src
		t = p
		return
	}
}
