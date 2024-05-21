package target

import (
	"sync"
)

func Cached[T Portal](finder Finder[T]) Finder[T] {
	store := newCache[T]()
	return func(resolve Path) Find[T] {
		resolveCached := func(src string) (path string, err error) {
			// try resolve from cache
			if portal, ok := store.Get(src); ok {
				return portal.Abs(), err
			}

			// try resolve from resolver
			path, err = resolve(src)
			return
		}

		find := finder(resolveCached)
		return func(src string) (portals Portals[T], err error) {
			portals, err = find(src)
			if err == nil {
				store.Add(portals)
			}
			return
		}
	}
}

type cache[T Portal] struct {
	portals Portals[T]
	mu      sync.Mutex
}

func newCache[T Portal]() *cache[T] {
	return &cache[T]{portals: make(Portals[T])}
}

func (c *cache[T]) Add(portals Portals[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for s, portal := range portals {
		c.portals[s] = portal
	}
}

func (c *cache[T]) Get(src string) (portal Portal, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range c.portals {
		m := p.Manifest()
		if m.Name == src || m.Package == src {
			ok = true
			portal = p
			return
		}
	}
	return
}
