package portals

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
	"log"
)

func Find(resolve target.Path) target.Find[target.Portal] {
	return resolver{Resolver: apps.Resolver{Path: resolve}}.resolve
}

type resolver struct{ apps.Resolver }

func (p resolver) resolve(src string) (portals target.Portals[target.Portal], err error) {
	portals = make(target.Portals[target.Portal])
	a, err1 := p.Resolver.Resolve(src)
	if err1 == nil {
		for s, app := range a {
			portals[s] = app
		}
	}

	projects, err2 := p.projects(src)
	if err2 == nil {
		for s, p := range projects {
			portals[s] = p
		}
	}
	for s, pp := range portals {
		log.Println("*", pp.Abs(), s)
	}
	if len(portals) > 0 {
		return
	}
	err = errors.Join(fmt.Errorf("cannot find portal %v for ", src), err1, err2)
	return
}

func (p resolver) projects(src string) (apps target.Portals[target.Project], err error) {
	apps = make(target.Portals[target.Project])
	for app := range sources.FromPath[target.Project](src) {
		if apps[app.Manifest().Package] == nil {
			apps[app.Manifest().Package] = app
		}
	}
	return
}
