package project

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
)

type NodeModule struct {
	*Module
	pkgJson bundle.PackageJson
}

func (m *Module) NodeModule() (module *NodeModule, err error) {
	pkgJson, err := bundle.LoadPackageJson(m.files)
	if err != nil {
		return
	}
	module = &NodeModule{Module: m, pkgJson: pkgJson}
	return
}

func (m *NodeModule) PkgJson() bundle.PackageJson {
	return m.pkgJson
}

func (m *NodeModule) IsPortalLib() bool {
	return m.pkgJson.IsPortalLib()
}

func (m *NodeModule) CanNpmRunBuild() bool {
	return m.pkgJson.Scripts.Build != ""
}

func (m *NodeModule) NpmRunBuild() (err error) {
	return exec.Run(m.src, "npm", "run", "build")
}

func (m *NodeModule) NpmInstall() (err error) {
	if err = exec.Run(m.src, "npm", "install"); err != nil {
		return fmt.Errorf("cannot install node_modules in %s: %s", m.src, err)
	}
	return
}

func (m *NodeModule) InjectDependencies(modules []NodeModule) (err error) {
	for _, module := range modules {
		if err = m.InjectDependency(module); err != nil {
			return
		}
	}
	return
}

func (m *NodeModule) InjectDependency(module NodeModule) (err error) {
	nm := path.Join(m.src, "node_modules", path.Base(module.Path()))
	log.Printf("copying module %v %v into: %s", module.Path(), module.pkgJson, nm)
	return fs.WalkDir(module.Files(), ".", func(s string, d fs.DirEntry, err error) error {
		path.Join(s, d.Name())
		if d.IsDir() {
			dst := path.Join(nm, s)
			if err = os.MkdirAll(dst, 0755); err != nil {
				return fmt.Errorf("os.MkdirAll: %v", err)
			}
			return nil
		}
		src, err := module.Files().Open(s)
		if err != nil {
			return fmt.Errorf("cannot open %s: %s", s, err)
		}
		defer src.Close()
		dst, err := os.Create(path.Join(nm, s))
		if err != nil {
			return fmt.Errorf("os.Create: %v", err)
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			return fmt.Errorf("io.Copy: %v", err)
		}
		return nil
	})
}
