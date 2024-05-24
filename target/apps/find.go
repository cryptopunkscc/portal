package apps

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
)

func NewFinder(
	getPath target.Path,
	files ...fs.FS,
) (f Finder) {
	f = Finder{
		GetPath: getPath,
		Files:   assets.ArrayFs(files),
	}
	return
}

type Finder struct {
	GetPath target.Path
	Files   fs.FS
}

func (a Finder) Find(ctx context.Context, src string) (apps target.Portals[target.App], err error) {
	log := plog.Get(ctx).Type(a).Set(&ctx)
	apps = make(target.Portals[target.App])
	tmp := src
	if src, _ = a.GetPath(src); src == "" {
		src = tmp
		log.Println("cannot resolve path for:", src)
	}

	if apps, err = a.ByPath(ctx, src); err != nil {
		log.Printf("cannot find apps for %s: %v", src, err)
		err = fmt.Errorf("apps.Finder cannot resolve app by path %v", src)
		return
	}
	return
}

func (a Finder) ByPath(ctx context.Context, src string) (apps target.Portals[target.App], err error) {
	log := plog.Get(ctx)

	apps = map[string]target.App{}
	if s := source.FromFS(a.Files, src).Lift(); s.Files() != nil {
		log.Println("Collecting from source", src)
		for _, app := range FromSource[target.App](s) {
			apps[app.Manifest().Package] = app
		}
	}

	log.Println("Collecting from path", src)
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
