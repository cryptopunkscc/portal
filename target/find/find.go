package find

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/source"
	"log"
)

type Deps[T target.Base] interface {
	TargetResolve() target.Resolve[T]
	TargetCache() *target.Cache[T]
	Path() target.Path
	Embed() []target.Source
	Priority() target.Priority
}

func Inject[T target.Base](deps Deps[T]) target.Find[T] {
	return (&finder[T]{
		resolveTarget: deps.TargetResolve(),
		cache:         deps.TargetCache(),
		resolvePath:   deps.Path(),
		embed:         deps.Embed(),
		priority:      deps.Priority(),
	}).Find
}

type finder[T target.Base] struct {
	cache         *target.Cache[T]
	resolvePath   target.Path
	resolveTarget target.Resolve[T]
	priority      target.Priority
	embed         []target.Source
}

func New[T target.Base](
	cache *target.Cache[T],
	resolvePath target.Path,
	resolveTarget target.Resolve[T],
	priority target.Priority,
	embed ...target.Source,
) target.Find[T] {
	return (&finder[T]{
		cache:         cache,
		resolvePath:   resolvePath,
		resolveTarget: resolveTarget,
		priority:      priority,
		embed:         embed,
	}).Find
}

func (f *finder[T]) Find(_ context.Context, src string) (apps target.Portals[T], err error) {
	if apps = f.fromEmbed(src); len(apps) > 0 {
		return
	}
	log.Println("getting from cache", src)
	if apps = f.fromCache(src); len(apps) > 0 {
		log.Println("got from cache", apps)
		return
	}
	log.Println("resolving path", src)
	if path, err := f.resolvePath(src); err == nil {
		log.Println("resoled path", path)
		src = path
	}
	log.Println("resolving from fs", src)
	if apps, err = f.fromFS(src); err != nil {
		err = ErrNoPortals
	}
	if len(apps) > 0 {
		apps.SortBy(f.priority)
		apps = apps.Reduced()
	}
	return
}

var ErrNoPortals = errors.New("cannot find portals")

func (f *finder[T]) fromEmbed(src string) (apps []T) {
	for _, t := range target.List(f.resolveTarget, f.embed...) {
		if t.Manifest().Match(src) {
			apps = append(apps, t)
		}
	}
	return
}

func (f *finder[T]) fromCache(src string) (apps []T) {
	if t, ok := f.cache.Get(src); ok {
		apps = append(apps, t)
	}
	return
}

func (f *finder[T]) fromFS(src string) (apps []T, err error) {
	file, err := source.File(src)
	if err == nil {
		apps = target.List(f.resolveTarget, file)
		f.cache.Add(apps)
	}
	return
}

func (f *finder[T]) getPriority(app T) int {
	for i, match := range f.priority {
		if match(app) {
			return i
		}
	}
	return len(f.priority)
}
