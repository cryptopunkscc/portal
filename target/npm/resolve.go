package npm

import (
	"encoding/json"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
	"io/fs"
)

var ErrNotNodeModule = errors.New("not a node module")

func ResolveNodeModule(m target.Source) (module target.NodeModule, err error) {
	if m.IsFile() {
		return nil, ErrNotNodeModule
	}
	m = m.Lift()
	pkgJson, err := loadPackageJson(m.Files())
	if err != nil {
		return
	}
	module = &source{Source: m, pkgJson: &pkgJson}
	return
}

func loadPackageJson(files fs.FS) (pkgJson target.PackageJson, err error) {
	file, err := fs.ReadFile(files, target.PackageJsonFilename)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &pkgJson)
	return
}
