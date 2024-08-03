package target

import "context"

func (f Find[T]) SortedBy(priority Priority) Find[T] {
	return func(ctx context.Context, src string) (portals Portals[T], err error) {
		if portals, err = f(ctx, src); err == nil {
			portals.SortBy(priority)
		}
		return
	}
}

func (f Find[T]) Reduced() Find[T] {
	return func(ctx context.Context, src string) (reduced Portals[T], err error) {
		if reduced, err = f(ctx, src); err != nil {
			reduced = reduced.Reduced()
		}
		return
	}
}
