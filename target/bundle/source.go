package bundle

import (
	"github.com/cryptopunkscc/portal/target"
)

type source struct {
	target.Source
	manifest *target.Manifest
}

type frontend struct {
	target.Html
	target.Bundle
}

type backend struct {
	target.Js
	target.Bundle
}

type executable struct {
	target.Exec
	target.Bundle
}

var _ target.BundleExec = &executable{}

func (b *source) IsApp() {}

func (b *source) IsBundle() {}

func (b *source) Manifest() *target.Manifest {
	return b.manifest
}
