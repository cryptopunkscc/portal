package portal

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
)

type of[T any] struct{ target.Portal_ }

func (a *of[T]) IsApp()        {}
func (a *of[T]) Target() (t T) { return }

func Resolve[T any](src target.Source) (t target.Portal[T], err error) {
	b, err := resolve(src)
	if err != nil {
		err = target.ErrNotTarget
		return
	}
	t = &of[T]{b}
	return
}

type of_ struct {
	target.Source
	manifest target.Manifest
}

func (p *of_) Manifest() *target.Manifest {
	return &p.manifest
}
func (p *of_) LoadManifest() error {
	return all.Unmarshalers.Load(&p.manifest, p.FS(), target.ManifestFilename)
}
func (p *of_) MarshalJSON() ([]byte, error) {
	return p.manifest.MarshalJSON()
}

func resolve(src target.Source) (t target.Portal_, err error) {
	t, ok := src.(target.Portal_)
	if ok {
		return
	}
	p := of_{Source: src}
	if err = p.LoadManifest(); err != nil {
		return
	}
	t = &p
	return
}

var Resolve_ target.Resolve[target.Portal_] = resolve
