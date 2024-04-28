package project

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	fs2 "github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"log"
	"path"
)

type NodeModule struct {
	Module
	pkgJson PackageJson
}

func (p *Module) NodeModule() (module *NodeModule, err error) {
	pkgJson, err := LoadPackageJson(p.files)
	if err != nil {
		return
	}
	module = &NodeModule{Module: *p, pkgJson: pkgJson}
	return
}

func (m *NodeModule) IsPortalModule() bool {
	return m.pkgJson.IsPortalModule()
}

func (m *NodeModule) HasNpmRunBuild() bool {
	return m.pkgJson.Scripts.Build != ""
}

func (m *NodeModule) NpmRunBuild() (err error) {
	return exec.Run(m.dir, "npm", "run", "build")
}

func (m *NodeModule) NpmInstall() (err error) {
	if err = exec.Run(m.dir, "npm", "install"); err != nil {
		return fmt.Errorf("cannot install node_modules: %s", err)
	}
	return
}

func (m *NodeModule) CopyModules(modules []string) (err error) {
	nodeModules := path.Join(m.dir, "node_modules")
	log.Printf("copying modules %v into: %s", modules, nodeModules)
	for _, module := range modules {
		dst := path.Join(nodeModules, path.Base(module))
		if err = fs2.CopyDir(module, dst); err != nil {
			return
		}
	}
	return
}
