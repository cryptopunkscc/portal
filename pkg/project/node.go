package project

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

type NodeModule struct {
	target.Source
	pkgJson *target.PackageJson
}

var _ target.NodeModule = (*NodeModule)(nil)

var ErrNotNodeModule = errors.New("not a node module")

func ResolveNodeModule(m target.Source) (module *NodeModule, err error) {
	if m.IsFile() {
		return nil, ErrNotNodeModule
	}
	m = m.Lift()
	pkgJson, err := target.LoadPackageJson(m.Files())
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
