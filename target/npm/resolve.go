package npm

import (
	"encoding/json"
	"errors"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
)

var ErrNotNodeModule = errors.New("not a node module")

func ResolveNodeModule(src target.Source) (nodeModule target.NodeModule, err error) {
	if src.IsFile() {
		return nil, ErrNotNodeModule
	}
	src = src.Lift()
	pkgJson, err := loadPackageJson(src.Files())
	if err != nil {
		return
	}
	nodeModule = &source{Source: src, pkgJson: &pkgJson}
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
