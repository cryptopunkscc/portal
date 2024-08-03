package find

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/target"
)

type Deps[T target.Portal_] interface {
	TargetResolve() target.Resolve[T]
	TargetCache() *target.Cache[T]
	TargetFile() target.File
	Path() target.Path
	Embed() []target.Source
	Priority() target.Priority
}

func Inject[T target.Portal_](deps Deps[T]) target.Find[T] {
	return New[T](
		deps.TargetCache(),
		deps.Path(),
		deps.TargetFile(),
		deps.TargetResolve(),
		deps.Priority(),
		deps.Embed()...,
	)
}

func New[T target.Portal_](
	cache *target.Cache[T],
	resolvePath target.Path,
	resolveFile target.File,
	resolveTarget target.Resolve[T],
	priority target.Priority,
	embed ...target.Source,
) target.Find[T] {
	return (&finder[T]{
		cache:         cache,
		resolvePath:   resolvePath,
		resolveFile:   resolveFile,
		resolveTarget: resolveTarget,
		priority:      priority,
		embed:         embed,
	}).Find
}

type finder[T target.Portal_] struct {
	cache         *target.Cache[T]
	resolvePath   target.Path
	resolveFile   target.File
	resolveTarget target.Resolve[T]
	priority      target.Priority
	embed         []target.Source
}

func (f *finder[T]) Find(_ context.Context, src string) (apps target.Portals[T], err error) {
	if apps = f.fromEmbed(src); len(apps) > 0 {
		return
	}
	if apps = f.fromCache(src); len(apps) > 0 {
		return
	}
	if path, err := f.resolvePath(src); err == nil {
		src = path
	}
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
	file, err := f.resolveFile(src)
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
