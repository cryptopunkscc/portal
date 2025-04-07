package target

import (
	"fmt"
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

func (c *Cache[T]) Add(portals ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, portal := range portals {
		c.portals()[portal.Manifest().Package] = portal
	}
}
func (c *Cache[T]) Get(src string) (portal T, ok bool) {
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

func (c *Cache[T]) Path(src string) (path string, err error) {
	if portal, ok := c.Get(src); ok {
		path = portal.Abs()
	} else {
		err = fmt.Errorf("no entry for %v", src)
	}
	return
}
