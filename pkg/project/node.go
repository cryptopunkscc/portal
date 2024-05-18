package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"path"
)

type NodeModule struct {
	target.Source
	pkgJson *target.PackageJson
}

var _ target.NodeModule = (*NodeModule)(nil)

func SkipNodeModulesDir(source target.Source) (result target.Source, err error) {
	if path.Base(source.Path()) == "node_modules" {
		return nil, fs.SkipDir
	}
	return
}

func ResolveNodeModule(m target.Source) (module *NodeModule, err error) {
	sub, err := fs.Sub(m.Files(), m.Path())
	if err != nil {
		return
	}
	pkgJson, err := target.LoadPackageJson(sub)
	if err != nil {
		return
	}
	module = &NodeModule{Source: m, pkgJson: &pkgJson}
	return
}

func (m *NodeModule) PkgJson() *target.PackageJson {
	return m.pkgJson
}

func (m *NodeModule) IsPortalLib() bool {
	return m.pkgJson.IsPortalLib()
}

func (m *NodeModule) CanNpmRunBuild() bool {
	return m.pkgJson.Scripts.Build != ""
}
