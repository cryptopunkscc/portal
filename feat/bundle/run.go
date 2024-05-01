package bundle

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/zip"
	"os"
	"path"
)

func RunAll(dir string) (err error) {
	root := path.Clean(dir)
	dir = "."
	if !path.IsAbs(dir) {
		dir = root
		root, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	found := false
	for app := range project.Find[project.PortalRawModule](os.DirFS(root), dir) {
		if err = Run(app); err != nil {
			return fmt.Errorf("bundle target %v: %v", app.Path(), err)
		}
		found = true
	}
	if !found {
		err = errors.New("no targets found")
	}
	return
}

func Run(app project.PortalRawModule) (err error) {

	// create build dir
	buildDir := path.Join(app.Parent().Path(), "build")
	if err = os.MkdirAll(buildDir, 0775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("os.MkdirAll: %v", err)
	}

	// pack dist dir
	bundleName := fmt.Sprintf("%s_%s.portal", app.Manifest().Name, app.Manifest().Version)
	if err = zip.Pack(app.Path(), path.Join(buildDir, bundleName)); err != nil {
		return fmt.Errorf("Pack: %v", err)
	}

	return
}

type PackageJson struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func (pkg *PackageJson) Load(src string) (err error) {
	bytes, err := os.ReadFile(src)
	if err != nil {
		return
	}
	return json.Unmarshal(bytes, pkg)
}
