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
	pkgJson *bundle.PackageJson
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
	module = &NodeModule{Source: m, pkgJson: &pkgJson}
	return
}

func (m *NodeModule) PkgJson() *bundle.PackageJson {
	return m.pkgJson
}

func (m *NodeModule) IsPortalLib() bool {
	return m.pkgJson.IsPortalLib()
}

func (m *NodeModule) CanNpmRunBuild() bool {
	return m.pkgJson.Scripts.Build != ""
}

func NpmRunBuild(m target.NodeModule) (err error) {
	return exec.Run(m.Abs(), "npm", "run", "build")
}

func NpmInstall(m target.NodeModule) (err error) {
	if err = exec.Run(m.Abs(), "npm", "install"); err != nil {
		return fmt.Errorf("cannot install node_modules in %s: %s", m.Abs(), err)
	}
	return
}

func InjectDependencies(m target.NodeModule, deps []target.NodeModule) (err error) {
	for _, module := range deps {
		if err = InjectDependency(m, module); err != nil {
			return
		}
	}
	return
}

func InjectDependency(m target.NodeModule, dep target.NodeModule) (err error) {
	nm := path.Join(m.Abs(), "node_modules", path.Base(dep.Abs()))
	log.Printf("copying module %v %v into: %s", dep.Abs(), dep.PkgJson(), nm)
	return fs.WalkDir(dep.Files(), ".", func(s string, d fs.DirEntry, err error) error {
		path.Join(s, d.Name())
		if d.IsDir() {
			dst := path.Join(nm, s)
			if err = os.MkdirAll(dst, 0755); err != nil {
				return fmt.Errorf("os.MkdirAll: %v", err)
			}
			return nil
		}
		src, err := dep.Files().Open(s)
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
