package portal

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
)

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
	for s, portal := range portals {
		log.Println("*", portal.Abs(), s)
	}
	if len(portals) > 0 {
		return
	}
	err = errors.Join(fmt.Errorf("cannot find portal %v for ", src), err1, err2)
	return
}

func ResolveProjects(src string) (apps target.Portals[target.Project], err error) {
	apps = make(target.Portals[target.Project])
	for app := range project.FindInPath[target.Project](src) {
		if apps[app.Manifest().Package] == nil {
			apps[app.Manifest().Package] = app
		}
	}
	return
}
