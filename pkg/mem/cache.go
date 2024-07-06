package mem

import (
	"maps"
	"reflect"
	"sync"
)

type Cache[T any] struct {
	entries  map[string]T
	mutex    sync.RWMutex
	onChange func(string, T, bool)
}

func (c *Cache[T]) Size() int {
	return len(c.entries)
}

func NewCache[T any]() Cache[T] {
	return Cache[T]{entries: make(map[string]T)}
}

func (c *Cache[T]) OnChange(onChange func(string, T, bool)) *Cache[T] {
	c.onChange = onChange
	return c
}

func (c *Cache[T]) Release() (entries map[string]T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entries = c.entries
	c.entries = map[string]T{}
	return
}

func (c *Cache[T]) Set(id string, t T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.SetUnsafe(id, t)
}

func (c *Cache[T]) SetUnsafe(id string, t T) {
	add := !reflect.ValueOf(t).IsZero()
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

func (c *Cache[T]) Delete(id string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.DeleteUnsafe(id)
}

func (c *Cache[T]) DeleteUnsafe(id string) (ok bool) {
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

func (c *Cache[T]) Get(id string) (t T, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	t, ok = c.entries[id]
	return
}

func (c *Cache[T]) Copy() map[string]T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return maps.Clone(c.entries)
}
