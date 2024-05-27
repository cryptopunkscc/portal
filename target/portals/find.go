package portals

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
	"io/fs"
	"strings"
)

var Finder target.Finder[target.Portal] = NewFind

func NewFind(getPath target.Path, files ...fs.FS) target.Find[target.Portal] {
	return finder{apps.NewFinder(getPath, files...)}.find
}

type finder struct{ apps.Finder }

func (p finder) find(ctx context.Context, src string) (portals target.Portals[target.Portal], err error) {
	base := src
	src = strings.TrimPrefix(src, "dev.")
	portals = make(target.Portals[target.Portal])

	if s, _ := p.GetPath(src); s != "" {
		src = s
	}

	if a, err := p.Finder.ByPath(ctx, src); err == nil {
		for s, app := range a {
			portals[s] = app
		}
	}

	for _, a := range sources.FromPath[target.Project](src) {
		if portals[a.Manifest().Package] == nil {
			portals[a.Manifest().Package] = a
		}
	}

	if len(portals) > 0 {
		return
	}
	err = fmt.Errorf("cannot find portal for %v", base)
	return
}
