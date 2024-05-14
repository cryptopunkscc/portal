package project

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/zip"
	"os"
	"path"
)

func BundlePortalApps(base, sub string) (err error) {
	found := false
	for app := range Find[PortalRawModule](os.DirFS(base), sub) {
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

func BundlePortalApp(app PortalRawModule) (err error) {

	// create build dir
	buildDir := path.Join(app.Parent().Path(), "build")
	if err = os.MkdirAll(buildDir, 0775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("os.MkdirAll: %v", err)
	}

	// pack dist dir
	bundleName := fmt.Sprintf("%s_%s.portal", app.Manifest().Name, app.Manifest().Version)
	if err = zip.Pack(app.Path(), path.Join(buildDir, bundleName)); err != nil {
		return fmt.Errorf("Pack: %v", err)
	}

	return
}
