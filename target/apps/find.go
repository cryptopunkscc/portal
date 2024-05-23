package apps

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/bundle"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
	"log"
	"strings"
)

func NewFinder(
	resolve target.Path,
	files ...fs.FS,
) Finder {
	return Finder{
		GetPath: resolve,
		Files:   assets.ArrayFs(files),
	}
}

type Finder struct {
	GetPath target.Path
	Files   fs.FS
}

func (a Finder) Find(src string) (apps target.Portals[target.App], err error) {
	log.Println("resolving app from:", src)
	apps = make(target.Portals[target.App])
	tmp := src
	if src, _ = a.GetPath(src); src == "" {
		src = tmp
	}
	log.Println("resolving app path:", src)

	if apps, err = a.ByPath(src); err != nil {
		err = fmt.Errorf("apps.Finder cannot resolve app by path %v", src)
		return
	}
	log.Println("resolved apps from:", src, apps)
	return
}

func (a Finder) ByNameOrPackage(src string) (app target.App, err error) {
	src = strings.TrimPrefix(src, "dev.")
	if src, err = a.GetPath(src); err != nil {
		return
	}
	var t target.Bundle
	if t, err = bundle.FromPath(src); err != nil {
		return nil, fmt.Errorf("cannot resolve portal apps path from '%s': %v", src, err)
	}
	app = t
	return
}

func (a Finder) ByPath(src string) (apps target.Portals[target.App], err error) {
	//log.Println("Finder.ByPath:", src)
	s := source.FromFS(a.Files, src)
	s = s.Lift()
	if s.Files() == nil {
		return
	}

	apps = map[string]target.App{}
	for _, app := range FromSource[target.App](s) {
		apps[app.Manifest().Package] = app
	}
	for _, app := range FromPath[target.App](src) {
		apps[app.Manifest().Package] = app
	}
	return
}

func FromPath[T target.App](src string) []T {
	return FromSource[T](source.FromPath(src))
}

func FromFS[T target.App](src fs.FS) []T {
	return FromSource[T](source.FromFS(src))
}

func FromSource[T target.App](src target.Source) []T {
	return source.List[T](Resolve[T](), src)
}

func Resolve[T target.App]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Try(dist.Resolve),
	)
}
