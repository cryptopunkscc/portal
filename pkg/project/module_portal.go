package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"os"
	"path"
)

var _ target.Project = &PortalNodeModule{}

type PortalNodeModule struct {
	*NodeModule
	manifest bundle.Manifest
}

func (m *NodeModule) PortalNodeModule() (module *PortalNodeModule, err error) {
	manifest := bundle.Manifest{}
	if err = manifest.LoadFs(m.files, "package.json"); err != nil {
		return
	}
	if err = manifest.LoadFs(m.files, bundle.PortalJson); err != nil {
		return
	}
	module = &PortalNodeModule{NodeModule: m, manifest: manifest}
	return
}

func (m *Module) PortalNodeModule() (module *PortalNodeModule, err error) {
	nm, err := m.NodeModule()
	if err != nil {
		return
	}
	return nm.PortalNodeModule()
}

func (m PortalNodeModule) Type() target.Type {
	return m.NodeModule.Type() + target.Dev
}

func (m *PortalNodeModule) Manifest() bundle.Manifest {
	return m.manifest
}

func (m *PortalNodeModule) PrepareBuild(dependencies ...NodeModule) (err error) {
	if err = m.Prepare(dependencies...); err != nil {
		return
	}
	if err = m.Build(); err != nil {
		return
	}
	return
}

func (m *PortalNodeModule) Prepare(dependencies ...NodeModule) (err error) {
	if err = m.NpmInstall(); err != nil {
		return
	}
	if err = m.InjectDependencies(dependencies); err != nil {
		return
	}
	return
}

func (m *PortalNodeModule) Build() (err error) {
	if !m.CanNpmRunBuild() {
		return errors.New("missing npm build in package.json")
	}
	if err = m.NpmRunBuild(); err != nil {
		return
	}
	if err = m.CopyIcon(); err != nil {
		return
	}
	if err = m.CopyManifest(); err != nil {
		return
	}
	return
}

func (m *PortalNodeModule) CopyIcon() (err error) {
	if m.manifest.Icon == "" {
		return
	}
	iconSrc := path.Join(m.src, m.manifest.Icon)
	iconName := "icon" + path.Ext(m.manifest.Icon)
	iconDst := path.Join(m.src, "dist", iconName)
	if err = fs.CopyFile(iconSrc, iconDst); err != nil {
		return
	}
	m.manifest.Icon = iconName
	return
}

func (m *PortalNodeModule) CopyManifest() (err error) {
	bytes, err := json.Marshal(m.manifest)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(m.src, "dist", bundle.PortalJson), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
