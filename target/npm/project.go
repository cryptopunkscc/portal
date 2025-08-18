package npm

import (
	"github.com/cryptopunkscc/portal/api/target"
)

type Project[T any] struct {
	target.Project[T]
	nodeModule target.NodeModule
}

func (p *Project[T]) PkgJson() *target.PackageJson { return p.nodeModule.PkgJson() }
func (p *Project[T]) Changed() bool {
	return target.Changed(p, "node_modules")
}
