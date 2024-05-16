package project

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
)

type NodeModule struct {
	target.Source
	pkgJson bundle.PackageJson
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
	module = &NodeModule{Source: m, pkgJson: pkgJson}
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
	return exec.Run(m.Abs(), "npm", "run", "build")
}

func (m *NodeModule) NpmInstall() (err error) {
	if err = exec.Run(m.Abs(), "npm", "install"); err != nil {
		return fmt.Errorf("cannot install node_modules in %s: %s", m.Abs(), err)
	}
	return
}

func (m *NodeModule) InjectDependencies(modules []target.NodeModule) (err error) {
	for _, module := range modules {
		if err = m.InjectDependency(module); err != nil {
			return
		}
	}
	return
}

func (m *NodeModule) InjectDependency(module target.NodeModule) (err error) {
	nm := path.Join(m.Abs(), "node_modules", path.Base(module.Abs()))
	log.Printf("copying module %v %v into: %s", module.Abs(), module.PkgJson(), nm)
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
