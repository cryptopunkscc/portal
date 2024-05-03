package portal

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"os"
)

type Apps map[string]target.App

func ResolveApps(src string) (apps Apps, err error) {
	apps = map[string]target.App{}

	if fs.Exists(src) {
		// scan src as path for portal apps
		apps, err = ResolveAppsByPath(src)
	} else {
		// resolve app path from appstore using given src as package name
		apps[src], err = ResolveAppByPackageName(src)
	}
	return
}

func ResolveAppsByPath(src string) (apps Apps, err error) {
	apps = map[string]target.App{}
	var base, sub string
	base, sub, err = project.Path(src)
	if err != nil {
		return nil, fmt.Errorf("cannot portal apps path: %v", err)
	}
	for app := range project.Find[target.App](os.DirFS(base), sub) {
		apps[app.Manifest().Package] = app
	}
	return
}

func ResolveAppByPackageName(src string) (app target.App, err error) {
	if src, err = appstore.Path(src); err != nil {
		return
	}
	var bundle *project.Bundle
	if bundle, err = project.NewModule(src).Bundle(); err != nil {
		return nil, fmt.Errorf("cannot portal apps path: %v", err)
	}
	app = bundle
	return
}
