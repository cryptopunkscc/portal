package target

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
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
				logger := plog.Get(ctx)
				for _, t := range portals {
					logger.Printf("added to cache: %v %v", t.Manifest(), t.Abs())
				}
			}
			return
		}
	}
}

type Cache[T Portal] struct {
	_portals Portals[T]
	mu       sync.Mutex
}

func (c *Cache[T]) portals() Portals[T] {
	if c._portals == nil {
		c._portals = make(Portals[T])
	}
	return c._portals
}
func (c *Cache[T]) Add(portals Portals[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for s, portal := range portals {
		c.portals()[s] = portal
	}
}

func (c *Cache[T]) Get(src string) (portal Portal, ok bool) {
	defer func() {
		log.Println("get from cache:", src, ok, portal, c.portals())
	}()
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range c.portals() {
		if p.Manifest().Match(src) {
			ok = true
			portal = p
			return
		}
	}
	return
}
