package npm

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

type source struct {
	target.Source
	pkgJson *target.PackageJson
}

func (m *source) PkgJson() *target.PackageJson {
	return m.pkgJson
}
