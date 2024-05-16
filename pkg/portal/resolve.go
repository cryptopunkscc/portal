package portal

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"path"
)

func ResolveApps(src string) (apps target.Portals[target.App], err error) {
	apps = make(target.Portals[target.App])
	if fs.Exists(src) {
		if path.Base(src) == src {
			if apps[src], err = ResolveAppByNameOrPackage(src); err == nil {
				return
			}
		}

		// scan src as path for portal apps
		if apps, err = ResolveAppsByPath(src); err == nil {
			return
		}

	} else {
		// resolve app path from appstore using given src as package name
		if apps[src], err = ResolveAppByNameOrPackage(src); err == nil {
			return
		}
	}

	apps = nil
	return
}

func ResolveAppByNameOrPackage(src string) (app target.App, err error) {
	if src, err = appstore.Path(src); err != nil {
		return
	}
	var bundle *project.Bundle
	if bundle, err = project.NewBundle(src); err != nil {
		return nil, fmt.Errorf("cannot resolve portal apps path from '%s': %v", src, err)
	}
	app = bundle
	return
}

func ResolveAppsByPath(src string) (apps target.Portals[target.App], err error) {
	apps = map[string]target.App{}
	for app := range project.FindInPath[target.App](src) {
		apps[app.Manifest().Package] = app
	}
	return
}
