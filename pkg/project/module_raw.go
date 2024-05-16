package project

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/zip"
	"io/fs"
	"os"
	"path"
)

type PortalRawModule struct {
	target.Source
	manifest bundle.Manifest
}

func ResolvePortalRawModule(m target.Source) (module *PortalRawModule, err error) {
	sub, err := fs.Sub(m.Files(), m.Path())
	if err != nil {
		return
	}
	manifest, err := bundle.ReadManifestFs(sub)
	if err != nil {
		return
	}
	module = &PortalRawModule{Source: m, manifest: manifest}
	return
}

func (m *PortalRawModule) App() {}

func (m PortalRawModule) Type() target.Type {
	return m.Source.Type() + target.Dev
}

func (m *PortalRawModule) Manifest() bundle.Manifest {
	return m.manifest
}

func BundlePortalApps(base, sub string) (err error) {
	found := false
	for app := range FindInPath[*PortalRawModule](path.Join(base, sub)) {
		if err = BundlePortalApp(app); err != nil {
			return fmt.Errorf("bundle target %v: %v", app.Path(), err)
		}
		found = true
	}
	if !found {
		err = errors.New("no targets found")
	}
	return
}

func BundlePortalApp(app *PortalRawModule) (err error) {
	// create build dir
	buildDir := path.Join(app.Parent().Abs(), "build")
	if err = os.MkdirAll(buildDir, 0775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("os.MkdirAll: %v", err)
	}

	// pack dist dir
	bundleName := fmt.Sprintf("%s_%s.portal", app.Manifest().Name, app.Manifest().Version)
	if err = zip.Pack(app.Abs(), path.Join(buildDir, bundleName)); err != nil {
		return fmt.Errorf("Pack: %v", err)
	}

	return
}
