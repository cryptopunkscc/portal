package find

import (
	"context"
	"errors"
	t "github.com/cryptopunkscc/portal/target"
)

func ByPath[T t.Portal_](file t.File, resolve t.Resolve[T]) t.Find[T] {
	return func(ctx context.Context, src string) (portals t.Portals[T], err error) {
		f, err := file(src)
		if err == nil {
			portals = resolve.List(f)
		}
		return
	}
}

func ById[T t.Portal_](resolve t.Resolve[T], sources ...t.Source) t.Find[T] {
	return func(ctx context.Context, src string) (portals t.Portals[T], err error) {
		for _, next := range resolve.List(sources...) {
			if next.Manifest().Match(src) {
				portals = append(portals, next)
			}
		}
		if len(portals) == 0 {
			err = ErrNothing
		}
		return
	}
}

func Combine[T t.Portal_](of ...t.Find[T]) t.Find[T] {
	return func(ctx context.Context, src string) (portals t.Portals[T], err error) {
		for _, find := range of {
			if portals, err = find(ctx, src); err == nil {
				return
			}
		}
		err = ErrNothing
		return
	}
}

var ErrNothing = errors.New("found nothing")
