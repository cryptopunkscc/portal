package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

type NodeModule struct {
	target.Source
	pkgJson *bundle.PackageJson
}

func ResolveNodeModule(m target.Source) (module *NodeModule, err error) {
	sub, err := fs.Sub(m.Files(), m.Path())
	if err != nil {
		return
	}
	pkgJson, err := bundle.LoadPackageJson(sub)
	if err != nil {
		return
	}
	module = &NodeModule{Source: m, pkgJson: &pkgJson}
	return
}

func (m *NodeModule) PkgJson() *bundle.PackageJson {
	return m.pkgJson
}

func (m *NodeModule) IsPortalLib() bool {
	return m.pkgJson.IsPortalLib()
}

func (m *NodeModule) CanNpmRunBuild() bool {
	return m.pkgJson.Scripts.Build != ""
}
