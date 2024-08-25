package target

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"sync"
)

type portalMap[T Portal_] map[string]T

type Cache[T Portal_] struct {
	_portals portalMap[T]
	mu       sync.Mutex
}

func (c *Cache[T]) portals() portalMap[T] {
	if c._portals == nil {
		c._portals = make(portalMap[T])
	}
	return c._portals
}

func (c *Cache[T]) Add(portals []T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, portal := range portals {
		c.portals()[portal.Manifest().Package] = portal
	}
}
func (c *Cache[T]) Get(ctx context.Context, src string) (portal T, ok bool) {
	defer func() {
		plog.Get(ctx).Type(c).Println("get from cache:", src, ok, portal, c.portals())
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
