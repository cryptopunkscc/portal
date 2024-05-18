package project

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"log"
)

func Resolve(resolve target.Path) func(src string) (portals target.Portals[target.Portal], err error) {
	return portals{Apps: portal.Apps{Path: resolve}}.Resolve
}

type portals struct{ portal.Apps }

func (p portals) Resolve(src string) (portals target.Portals[target.Portal], err error) {
	portals = make(target.Portals[target.Portal])
	apps, err1 := p.Apps.Resolve(src)
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
	for s, pp := range portals {
		log.Println("*", pp.Abs(), s)
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

func FromPath[T target.Source](src string) (in <-chan T) {
	return target.Stream[T](Dev[T](), target.NewModule(src))
}

func FromFS[T target.Source](src fs.FS) (in <-chan T) {
	return target.Stream[T](Dev[T](), target.NewModuleFS(src))
}

func Dev[T target.Source]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(portal.ResolveBundle),
		target.Lift(target.Try(ResolveNodeModule))(
			target.Try(ResolvePortalModule)),
		target.Try(portal.ResolveDist),
	)
}
