package project

import (
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	fs2 "github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"log"
	"os"
	"path"
)

type NodeModule struct {
	*Module
	pkgJson PackageJson
}

func (m *Module) NodeModule() (module *NodeModule, err error) {
	pkgJson, err := LoadPackageJson(m.files)
	if err != nil {
		return
	}
	module = &NodeModule{Module: m, pkgJson: pkgJson}
	return
}

func (m *NodeModule) IsPortalModule() bool {
	return m.pkgJson.IsPortalModule()
}

func (m *NodeModule) HasNpmRunBuild() bool {
	return m.pkgJson.Scripts.Build != ""
}

func (m *NodeModule) NpmRunBuild() (err error) {
	return exec.Run(m.src, "npm", "run", "build")
}

func (m *NodeModule) NpmInstall() (err error) {
	if err = exec.Run(m.src, "npm", "install"); err != nil {
		return fmt.Errorf("cannot install node_modules: %s", err)
	}
	return
}

func (m *NodeModule) CopyManifest() (err error) {
	src := m.src
	b := bundle.Base(src)
	_ = b.LoadPath(src, "package.json")
	_ = b.LoadPath(src, bundle.PortalJson)
	if b.Icon != "" {
		iconSrc := path.Join(src, b.Icon)
		iconName := "icon" + path.Ext(b.Icon)
		iconDst := path.Join(src, "dist", iconName)
		if err = fs2.CopyFile(iconSrc, iconDst); err != nil {
			return
		}
		b.Icon = iconName
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(src, "dist", bundle.PortalJson), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}

func (m *NodeModule) CopyModules(modules []string) (err error) {
	nodeModules := path.Join(m.src, "node_modules")
	log.Printf("copying modules %v into: %s", modules, nodeModules)
	for _, module := range modules {
		dst := path.Join(nodeModules, path.Base(module))
		if err = fs2.CopyDir(module, dst); err != nil {
			return
		}
	}
	return
}
