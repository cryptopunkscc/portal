package portals

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/sources"
	"io/fs"
	"strings"
)

var Finder target.Finder[target.Portal] = NewFind[target.Portal]

func NewFind[T target.Portal](getPath target.Path, files ...fs.FS) target.Find[T] {
	return finder[T]{apps.NewFinder[target.App](getPath, files...)}.find
}

type finder[T target.Portal] struct{ apps.Finder[target.App] }

func (p finder[T]) find(ctx context.Context, src string) (portals target.Portals[T], err error) {
	base := src
	src = strings.TrimPrefix(src, "dev.")
	portals = make(target.Portals[T])

	if s, _ := p.GetPath(src); s != "" {
		src = s
	}

	if a, err := p.Finder.ByPath(ctx, src); err == nil {
		for s, app := range a {
			if t, ok := app.(T); ok {
				portals[s] = t
			}
		}
	}

	for _, a := range sources.FromPath[T](src) {
		if _, ok := any(a).(target.Project); ok {
			portals[a.Manifest().Package] = a
		}
	}

	if len(portals) > 0 {
		return
	}
	err = fmt.Errorf("cannot find portal for %v", base)
	return
}
