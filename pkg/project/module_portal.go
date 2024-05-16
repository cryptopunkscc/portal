package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	fs2 "io/fs"
	"os"
	"path"
)

var _ target.Project = &PortalNodeModule{}

type PortalNodeModule struct {
	target.NodeModule
	manifest *bundle.Manifest
}

func NewPortalNodeModule(src string) (module *PortalNodeModule, err error) {
	nodeModule, err := ResolveNodeModule(NewModule(src))
	if err != nil {
		return
	}
	return ResolvePortalNodeModule(nodeModule)
}

func ResolvePortalNodeModule(m target.NodeModule) (module *PortalNodeModule, err error) {
	manifest := bundle.Manifest{}
	sub, err := fs2.Sub(m.Files(), m.Path())
	if err != nil {
		return
	}
	if err = manifest.LoadFs(sub, "package.json"); err != nil {
		return
	}
	if err = manifest.LoadFs(sub, bundle.PortalJson); err != nil {
		return
	}
	module = &PortalNodeModule{NodeModule: m, manifest: &manifest}
	return
}

func (m PortalNodeModule) Type() target.Type {
	return m.NodeModule.Type() + target.Dev
}

func (m *PortalNodeModule) Manifest() *bundle.Manifest {
	return m.manifest
}

func BuildPortalApps(root, dir string, dependencies ...target.NodeModule) (err error) {
	for m := range FindInPath[target.Project](path.Join(root, dir)) {

		if !m.CanNpmRunBuild() {
			continue
		}
		if err = PrepareBuild(m, dependencies); err != nil {
			return err
		}
	}
	return
}

func PrepareBuild(m target.Project, dependencies []target.NodeModule) (err error) {
	if err = Prepare(m, dependencies); err != nil {
		return
	}
	if err = Build(m); err != nil {
		return
	}
	return
}

func Prepare(m target.Project, dependencies []target.NodeModule) (err error) {
	if err = NpmInstall(m); err != nil {
		return
	}
	if err = InjectDependencies(m, dependencies); err != nil {
		return
	}
	return
}

func Build(m target.Project) (err error) {
	if !m.CanNpmRunBuild() {
		return errors.New("missing npm build in package.json")
	}
	if err = NpmRunBuild(m); err != nil {
		return
	}
	if err = CopyIcon(m); err != nil {
		return
	}
	if err = CopyManifest(m); err != nil {
		return
	}
	return
}

func CopyIcon(m target.Project) (err error) {
	if m.Manifest().Icon == "" {
		return
	}
	iconSrc := path.Join(m.Abs(), m.Manifest().Icon)
	iconName := "icon" + path.Ext(m.Manifest().Icon)
	iconDst := path.Join(m.Abs(), "dist", iconName)
	if err = fs.CopyFile(iconSrc, iconDst); err != nil {
		return
	}
	m.Manifest().Icon = iconName
	return
}

func CopyManifest(m target.Project) (err error) {
	bytes, err := json.Marshal(m.Manifest())
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(m.Abs(), "dist", bundle.PortalJson), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
