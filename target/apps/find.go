package apps

import (
	"fmt"
	fsUtil "github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/bundle"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
	"log"
	"strings"
)

func Find(resolve target.Path) func(src string) (apps target.Portals[target.App], err error) {
	return Finder{resolve}.Find
}

type Finder struct{ target.Path }

func (a Finder) Find(src string) (apps target.Portals[target.App], err error) {
	apps = make(target.Portals[target.App])
	if !fsUtil.Exists(src) {
		tmp := src
		if src, err = a.Path(src); err != nil {
			err = fmt.Errorf("apps.Finder cannot resolve path from %v: %v", tmp, err)
			return
		}
	}
	log.Println("resolving app from:", src)
	if apps, err = a.ByPath(src); err != nil {
		err = fmt.Errorf("apps.Finder cannot resolve app by path %v", src)
		return
	}
	return
}

func (a Finder) ByNameOrPackage(src string) (app target.App, err error) {
	src = strings.TrimPrefix(src, "dev.")
	if src, err = a.Path(src); err != nil {
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
	apps = map[string]target.App{}
	for app := range FromPath[target.App](src) {
		apps[app.Manifest().Package] = app
	}
	return
}

func FromPath[T target.App](src string) (in <-chan T) {
	return source.Stream[T](Resolve[T](), source.FromPath(src))
}

func FromFS[T target.App](src fs.FS) (in <-chan T) {
	return source.Stream[T](Resolve[T](), source.FromFS(src))
}

func Resolve[T target.App]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Try(dist.Resolve),
	)
}
