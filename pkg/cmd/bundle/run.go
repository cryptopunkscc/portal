package bundle

import (
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/cryptopunkscc/go-astral-js/pkg/zip"
	"io/fs"
	"os"
	"path"
)

func RunAll(dir string) (err error) {
	r, err := runner.New(dir, runner.DevTargets)
	if err != nil {
		return
	}

	targets := append(r.Backends, r.Frontends...)

	for _, target := range targets {
		if err = Run(target.Path); err != nil {
			return fmt.Errorf("bundle target %v: %v", target.Path, err)
		}
	}

	return
}

func Run(src string) (err error) {
	srcFs := os.DirFS(src)

	// build dist if needed
	if _, err = fs.Stat(srcFs, "dist"); os.IsNotExist(err) {
		if err = build.Run(src); err != nil {
			return
		}
	}

	dist := path.Join(src, "dist")

	// prepare portal.json
	portalJson := &PackageJson{
		Name:    path.Base(src),
		Version: "0.0.0",
	}
	_ = portalJson.Load(path.Join(src, "portal.json"))
	_ = portalJson.Load(path.Join(src, "package.json"))
	bytes, err := json.Marshal(portalJson)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(dist, "portal.json"), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}

	// create build dir
	buildDir := path.Join(src, "/", "build")
	if err = os.MkdirAll(buildDir, 0775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("os.MkdirAll: %v", err)
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
