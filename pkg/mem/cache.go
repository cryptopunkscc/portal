package mem

import (
	"maps"
	"sync"
)

type cache[T any] struct {
	entries  map[string]T
	mutex    sync.RWMutex
	onChange func(string, T, bool)
}

func (c *cache[T]) Size() int {
	return len(c.entries)
}

func NewCache[T any]() Cache[T] {
	return &cache[T]{entries: make(map[string]T)}
}

func (c *cache[T]) OnChange(onChange func(string, T, bool)) {
	c.onChange = onChange
}

func (c *cache[T]) Release() (entries map[string]T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entries = c.entries
	c.entries = map[string]T{}
	return
}

func (c *cache[T]) Set(id string, t T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.SetUnsafe(id, t)
}

func (c *cache[T]) SetUnsafe(id string, t T) {
	add := any(t) != nil
	ok := true
	if add {
		c.entries[id] = t
	} else {
		t, ok = c.entries[id]
		delete(c.entries, id)
	}
	if c.onChange != nil && ok {
		c.onChange(id, t, add)
	}
}

func (c *cache[T]) Delete(id string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.DeleteUnsafe(id)
}

func (c *cache[T]) DeleteUnsafe(id string) (ok bool) {
	t, ok := c.entries[id]
	if !ok {
		return
	}
	delete(c.entries, id)

	if c.onChange != nil {
		c.onChange(id, t, false)
	}
	return
}

func (c *cache[T]) Get(id string) (t T, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	t, ok = c.entries[id]
	return
}

func (c *cache[T]) Copy() map[string]T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return maps.Clone(c.entries)
}
