package portal

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"os"
)

type Resolve[T target.Portal] func(src string) (apps target.Portals[T], err error)

type Apps target.Portals[target.App]

func ResolveApps(src string) (apps target.Portals[target.App], err error) {
	apps = make(target.Portals[target.App])
	if fs.Exists(src) {
		// scan src as path for portal apps
		apps, err = ResolveAppsByPath(src)
	} else {
		// resolve app path from appstore using given src as package name
		apps[src], err = ResolveAppByPackageName(src)
		if err != nil {
			apps = nil
		}
	}
	return
}

func ResolveAppsByPath(src string) (apps target.Portals[target.App], err error) {
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

type Projects target.Portals[target.Project]

func ResolveProjects(src string) (apps target.Portals[target.Project], err error) {
	apps = make(target.Portals[target.Project])
	var base, sub string
	base, sub, err = project.Path(src)
	if err != nil {
		return nil, fmt.Errorf("cannot portal apps path: %v", err)
	}
	for app := range project.Find[target.Project](os.DirFS(base), sub) {
		if apps[app.Manifest().Package] == nil {
			apps[app.Manifest().Package] = app
		}
	}
	return
}

func ResolvePortals(src string) (portals target.Portals[target.Portal], err error) {
	portals = make(target.Portals[target.Portal])
	apps, err1 := ResolveApps(src)
	if err1 == nil {
		for s, app := range apps {
			portals[s] = app
		}
	}

	projects, err2 := ResolveProjects(src)
	if err2 == nil {
		for s, p := range projects {
			portals[s] = p
		}
	}
	if len(portals) > 0 {
		return
	}
	err = errors.Join(fmt.Errorf("cannot find portal %v for ", src), err1, err2)
	return
}
