package resolve

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"path"
	"strings"
)

func Apps(resolve target.Path) func(src string) (apps target.Portals[target.App], err error) {
	return apps{resolve}.Resolve
}

type apps struct{ target.Path }

func (a apps) Resolve(src string) (apps target.Portals[target.App], err error) {
	apps = make(target.Portals[target.App])
	if fs.Exists(src) {
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

func (a apps) ByNameOrPackage(src string) (app target.App, err error) {
	src = strings.TrimPrefix(src, "dev.")
	if src, err = a.Path(src); err != nil {
		return
	}
	var bundle *project.Bundle
	if bundle, err = project.NewBundle(src); err != nil {
		return nil, fmt.Errorf("cannot resolve portal apps path from '%s': %v", src, err)
	}
	app = bundle
	return
}

func (a apps) ByPath(src string) (apps target.Portals[target.App], err error) {
	apps = map[string]target.App{}
	for app := range FromPath[target.App](src) {
		apps[app.Manifest().Package] = app
	}
	return
}
