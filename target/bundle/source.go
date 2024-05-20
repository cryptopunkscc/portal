package bundle

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

type source struct {
	target.Source
	manifest *target.Manifest
}

type frontend struct {
	target.Frontend
	target.Bundle
}

type backend struct {
	target.Backend
	target.Bundle
}

func (b *source) IsApp() {}

func (b *source) IsBundle() {}

func (b *source) Manifest() *target.Manifest {
	return b.manifest
}
