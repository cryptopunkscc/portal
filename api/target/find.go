package target

import (
	"context"
	"slices"
)

type (
	Find[T Portal_] func(ctx context.Context, src string) (portals Portals[T], err error)
	File            func(path ...string) (source Source, err error)
	Path            func(src string) (path string, err error)
)

func FindByPath[T Portal_](file File, resolve Resolve[T]) Find[T] {
	return func(ctx context.Context, src string) (portals Portals[T], err error) {
		f, err := file(src)
		if err == nil {
			portals = resolve.List(f)
		}
		return
	}
}

func (find Find[T]) Call(ctx context.Context, src string) (portals Portals[T], err error) {
	return find(ctx, src)
}

func (find Find[T]) ById(path Path) Find[T] {
	return func(ctx context.Context, src string) (portals Portals[T], err error) {
		if resolved, err := path(src); err == nil {
			src = resolved
		}
		return find(ctx, src)
	}
}

func (find Find[T]) Reduced(priority ...Matcher) Find[T] {
	return func(ctx context.Context, src string) (portals Portals[T], err error) {
		if portals, err = find(ctx, src); err == nil {
			portals.SortBy(priority)
			portals = portals.Reduced()
		}
		return
	}
}

func (find Find[T]) Cached(cache *Cache[T]) Find[T] {
	return func(ctx context.Context, src string) (portals Portals[T], err error) {
		if t, ok := cache.Get(src); ok {
			portals = append(portals, t)
			return
		}
		if portals, err = find(ctx, src); err != nil {
			return
		}
		cache.Add(portals)
		return
	}
}

func (p Portals[T]) SortBy(priority Priority) {
	slices.SortFunc(p, func(a, b T) int {
		return priority.Get(a) - priority.Get(b)
	})
}

func (p Portals[T]) Reduced() (reduced Portals[T]) {
	mem := make(map[string]T)
	for _, t := range p {
		if _, ok := mem[t.Manifest().Package]; ok {
			continue
		}
		mem[t.Manifest().Package] = t
		reduced = append(reduced, t)
	}
	return
}
