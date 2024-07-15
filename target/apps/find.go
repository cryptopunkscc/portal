package apps

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/assets"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/source"
	"io/fs"
)

func NewFind[T target.App](getPath target.Path, files ...fs.FS) target.Find[T] {
	return NewFinder[T](getPath, files...).Find
}

func NewFinder[T target.App](
	getPath target.Path,
	files ...fs.FS,
) (f Finder[T]) {
	f = Finder[T]{
		GetPath: getPath,
		Files:   assets.ArrayFs(files),
	}
	return
}

type Finder[T target.App] struct {
	GetPath target.Path
	Files   fs.FS
}

func (a Finder[T]) Find(ctx context.Context, src string) (apps target.Portals[T], err error) {
	log := plog.Get(ctx).Type(a)
	apps = make(target.Portals[T])
	if p, _ := a.GetPath(src); p != "" {
		src = p
		log.Println("resolved path for:", src)
	}

	if apps, err = a.ByPath(ctx, src); err != nil {
		log.Printf("cannot find apps for %s: %v", src, err)
		err = fmt.Errorf("apps.Finder cannot resolve app by path %v", src)
		return
	}
	return
}

func (a Finder[T]) ByPath(ctx context.Context, src string) (apps target.Portals[T], err error) {
	log := plog.Get(ctx).Type(a)

	apps = map[string]T{}
	if s := source.FromFS(a.Files, src).Lift(); s.Files() != nil {
		log.Println("Collecting from source", src)
		for _, app := range FromSource[T](s) {
			apps[app.Manifest().Package] = app
		}
	}

	log.Println("Collecting from path", src)
	for _, app := range FromPath[T](src) {
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
