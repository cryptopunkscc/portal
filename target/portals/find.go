package portals

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
	"io/fs"
	"strings"
)

func Find(
	getPath target.Path,
	files ...fs.FS,
) target.Find[target.Portal] {
	return Finder{apps.NewFinder(getPath, files...)}.find
}

type Finder struct{ apps.Finder }

func (p Finder) find(src string) (portals target.Portals[target.Portal], err error) {
	base := src
	src = strings.TrimPrefix(src, "dev.")
	portals = make(target.Portals[target.Portal])

	if s, _ := p.GetPath(src); s != "" {
		src = s
	}

	if a, err1 := p.Finder.ByPath(src); err1 == nil {
		for s, app := range a {
			portals[s] = app
		}
	}

	if projects, err2 := p.projects(src); err2 == nil {
		for s, t := range projects {
			portals[s] = t
		}
	}

	if len(portals) > 0 {
		return
	}
	err = fmt.Errorf("cannot find portal for %v", base)
	return
}

func (p Finder) projects(src string) (apps target.Portals[target.Project], err error) {
	apps = make(target.Portals[target.Project])
	for _, app := range sources.FromPath[target.Project](src) {
		if apps[app.Manifest().Package] == nil {
			apps[app.Manifest().Package] = app
		}
	}
	return
}
