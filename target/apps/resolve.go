package apps

import (
	"fmt"
	fsUtil "github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/bundle"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
	"path"
	"strings"
)

func Resolve(resolve target.Path) func(src string) (apps target.Portals[target.App], err error) {
	return Resolver{resolve}.Resolve
}

type Resolver struct{ target.Path }

func (a Resolver) Resolve(src string) (apps target.Portals[target.App], err error) {
	apps = make(target.Portals[target.App])
	if fsUtil.Exists(src) {
		if path.Base(src) == src {
			if apps[src], err = a.ByNameOrPackage(src); err == nil {
				return
			}
		}

		if apps, err = a.ByPath(src); err == nil {
			return
		}

	} else {
		if apps[src], err = a.ByNameOrPackage(src); err == nil {
			return
		}
	}

	apps = nil
	return
}

func (a Resolver) ByNameOrPackage(src string) (app target.App, err error) {
	src = strings.TrimPrefix(src, "dev.")
	if src, err = a.Path(src); err != nil {
		return
	}
	var t target.Bundle
	if t, err = bundle.New(src); err != nil {
		return nil, fmt.Errorf("cannot resolve portal apps path from '%s': %v", src, err)
	}
	app = t
	return
}

func (a Resolver) ByPath(src string) (apps target.Portals[target.App], err error) {
	apps = map[string]target.App{}
	for app := range FromPath[target.App](src) {
		apps[app.Manifest().Package] = app
	}
	return
}

func FromPath[T target.App](src string) (in <-chan T) {
	return source.Stream[T](resolve[T](), source.New(src))
}

func FromFS[T target.App](src fs.FS) (in <-chan T) {
	return source.Stream[T](resolve[T](), source.Resolve(src))
}

func resolve[T target.App]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Try(dist.Resolve),
	)
}
