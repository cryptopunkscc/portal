package bundle

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/zip"
	"io/fs"
	"log"
	"os"
	"path"
)

func RunAll(dir string) (err error) {
	found := false
	for target := range project.ProdTargets(os.DirFS(dir)) {
		log.Println(target.Path())
		if err = Run(target.Path()); err != nil {
			return fmt.Errorf("bundle target %v: %v", target.Path(), err)
		}
		found = true
	}
	if !found {
		err = errors.New("no targets found")
	}
	return
}

func Run(src string) (err error) {
	srcFs := os.DirFS(src)

	// build dist if needed
	if stat, err := fs.Stat(srcFs, "package.json"); err == nil && stat.Mode().IsRegular() {
		if _, err = fs.Stat(srcFs, "dist"); os.IsNotExist(err) {
			if err = build.Run(src); err != nil {
				return err
			}
		}
	}

	// load manifest
	portalJson := bundle.Base(src)
	if err = portalJson.LoadPath(src, bundle.PortalJson); err != nil {
		return fmt.Errorf("portalJson.LoadPath: %v", err)
	}

	// create build dir
	buildDir := path.Join(src, "/", "build")
	if err = os.MkdirAll(buildDir, 0775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("os.MkdirAll: %v", err)
	}

	// resolve dist dir
	dist := src
	if stat, err := fs.Stat(srcFs, "dist"); err == nil && stat.IsDir() {
		dist = path.Join(src, "dist")
	}

	// pack dist dir
	bundleName := fmt.Sprintf("%s_%s.portal", portalJson.Name, portalJson.Version)
	if err = zip.Pack(dist, path.Join(buildDir, bundleName)); err != nil {
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
