package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"os"
	"path"
)

type PortalNodeModule struct {
	*NodeModule
	manifest bundle.Manifest
}

func (m *NodeModule) PortalNodeModule() (module *PortalNodeModule, err error) {
	manifest, err := bundle.ReadManifestFs(m.files)
	if err != nil {
		return
	}
	module = &PortalNodeModule{NodeModule: m, manifest: manifest}
	return
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
	if err = m.CopyManifest(); err != nil {
		return
	}
	return
}

func (m *PortalNodeModule) CopyManifest() (err error) {
	src := m.src
	b := bundle.Base(src)
	_ = b.LoadPath(src, "package.json")
	_ = b.LoadPath(src, bundle.PortalJson)
	if b.Icon != "" {
		iconSrc := path.Join(src, b.Icon)
		iconName := "icon" + path.Ext(b.Icon)
		iconDst := path.Join(src, "dist", iconName)
		if err = fs.CopyFile(iconSrc, iconDst); err != nil {
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
