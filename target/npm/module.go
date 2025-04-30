package npm

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/dec/json"
)

type Module struct {
	target.Source
	packageJson *target.PackageJson
}

func (n *Module) PkgJson() *target.PackageJson {
	return n.packageJson
}
func (n *Module) LoadPkgJson() error {
	return json.Unmarshaler.Load(&n.packageJson, n.FS(), target.PackageJsonFilename)
}
