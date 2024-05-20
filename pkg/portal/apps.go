package portal

import (
	"fmt"
	fsUtil "github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"path"
	"strings"
)

func Resolve(resolve target.Path) func(src string) (apps target.Portals[target.App], err error) {
	return Apps{resolve}.Resolve
}

type Apps struct{ target.Path }

func (a Apps) Resolve(src string) (apps target.Portals[target.App], err error) {
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

func (a Apps) ByNameOrPackage(src string) (app target.App, err error) {
	src = strings.TrimPrefix(src, "dev.")
	if src, err = a.Path(src); err != nil {
		return
	}
	var bundle target.Bundle
	if bundle, err = NewBundle(src); err != nil {
		return nil, fmt.Errorf("cannot resolve portal apps path from '%s': %v", src, err)
	}
	app = bundle
	return
}

func (a Apps) ByPath(src string) (apps target.Portals[target.App], err error) {
	apps = map[string]target.App{}
	for app := range FromPath[target.App](src) {
		apps[app.Manifest().Package] = app
	}
	return
}

func FromPath[T target.App](src string) (in <-chan T) {
	return target.Stream[T](App[T](), target.NewModule(src))
}

func FromFS[T target.App](src fs.FS) (in <-chan T) {
	return target.Stream[T](App[T](), target.NewModuleFS(src))
}

func App[T target.App]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(ResolveBundle),
		target.Try(ResolveDist),
	)
}
