package target

import (
	"context"
	"io/fs"
	"log"
	"sync"
)

func (finder Finder[T]) Cached(c *Cache[T]) Finder[T] {
	return func(resolve Path, files ...fs.FS) Find[T] {
		resolveCached := func(src string) (path string, err error) {
			// try resolve from cache
			if portal, ok := c.Get(src); ok {
				return portal.Abs(), err
			}

			// try resolve from resolver
			path, err = resolve(src)
			return
		}

		find := finder(resolveCached, files...)
		return func(ctx context.Context, src string) (portals Portals[T], err error) {
			portals, err = find(ctx, src)
			if err == nil {
				c.Add(portals)
			}
			return
		}
	}
}

type Cache[T Portal] struct {
	portals Portals[T]
	mu      sync.Mutex
}

func NewCache[T Portal]() *Cache[T] {
	return &Cache[T]{portals: make(Portals[T])}
}

func (c *Cache[T]) Add(portals Portals[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for s, portal := range portals {
		c.portals[s] = portal
	}
	log.Println("added to cache:", c.portals)
}

func (c *Cache[T]) Get(src string) (portal Portal, ok bool) {
	defer func() {
		log.Println("get from cache:", src, ok, portal, c.portals)
	}()
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
