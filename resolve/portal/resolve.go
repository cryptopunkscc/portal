package portal

import (
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/cryptopunkscc/portal/target"
)

type of[T any] struct{ target.Portal_ }

func (a *of[T]) IsApp()        {}
func (a *of[T]) Target() (t T) { return }

func Resolve[T any](src target.Source) (t target.Portal[T], err error) {
	b, err := resolve(src)
	if err != nil {
		return
	}
	t = &of[T]{b}
	return
}

type unknown struct {
	target.Source
	manifest target.Manifest
}

func (p *unknown) Manifest() *target.Manifest {
	return &p.manifest
}

func (p *unknown) LoadManifest() error {
	return all.Unmarshalers.Load(&p.manifest, p.Files(), target.ManifestFilename)
}

func resolve(src target.Source) (t target.Portal_, err error) {
	p := unknown{Source: src}
	if err = p.LoadManifest(); err != nil {
		return
	}
	t = &p
	return
}

var Resolve_ target.Resolve[target.Portal_] = resolve
