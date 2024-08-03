package npm

import (
	"github.com/cryptopunkscc/portal/pkg/json"
	"github.com/cryptopunkscc/portal/target"
)

type nodeModule struct {
	target.Source
	packageJson *target.PackageJson
}

func (n *nodeModule) PkgJson() *target.PackageJson {
	return n.packageJson
}

func (n *nodeModule) LoadPkgJson() error {
	return json.Load(&n.packageJson, n.Files(), target.PackageJsonFilename)
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
