package resolve

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
)

func Portals(resolve target.Path) func(src string) (portals target.Portals[target.Portal], err error) {
	return portals{resolve}.Resolve
}

type portals apps

func (p portals) Resolve(src string) (portals target.Portals[target.Portal], err error) {
	portals = make(target.Portals[target.Portal])
	apps, err1 := p.Resolve(src)
	if err1 == nil {
		for s, app := range apps {
			portals[s] = app
		}
	}

	projects, err2 := p.Projects(src)
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

func (p portals) Projects(src string) (apps target.Portals[target.Project], err error) {
	apps = make(target.Portals[target.Project])
	for app := range FromPath[target.Project](src) {
		if apps[app.Manifest().Package] == nil {
			apps[app.Manifest().Package] = app
		}
	}
	return
}
