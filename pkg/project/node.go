package project

import (
	"encoding/json"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
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
	pkgJson, err := loadPackageJson(m.Files())
	if err != nil {
		return
	}
	module = &NodeModule{Source: m, pkgJson: &pkgJson}
	return
}

func (m *NodeModule) PkgJson() *target.PackageJson {
	return m.pkgJson
}

func loadPackageJson(files fs.FS) (pkgJson target.PackageJson, err error) {
	file, err := fs.ReadFile(files, target.PackageJsonFilename)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &pkgJson)
	return
}
