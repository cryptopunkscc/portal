package npm

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/json"
)

type nodeModule struct {
	target2.Source
	packageJson *target.PackageJson
}

func (n *nodeModule) PkgJson() *target.PackageJson {
	return n.packageJson
}

func (n *nodeModule) LoadPkgJson() error {
	return json.Load(&n.packageJson, n, target.PackageJsonFilename)
}

func Resolve(src target2.Source) (t target2.NodeModule, err error) {
	s := &nodeModule{Source: src}
	if err = s.LoadPkgJson(); err != nil {
		return
	}
	t = s
	return
}
