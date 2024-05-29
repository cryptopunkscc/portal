package apphost

import (
	"reflect"
	"sync"
)

type Registry[T any] struct {
	entries  map[string]T
	mutex    sync.RWMutex
	onChange func(map[string]T, string, bool)
}

func (r *Registry[T]) Release() (entries map[string]T) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	entries = r.entries
	r.entries = map[string]T{}
	return
}

func (r *Registry[T]) Set(id string, t T) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.SetUnsafe(id, t)
}

func (r *Registry[T]) SetUnsafe(id string, t T) {
	add := !reflect.ValueOf(t).IsZero()
	if add {
		r.entries[id] = t
	} else {
		delete(r.entries, id)
	}
	if r.onChange != nil {
		r.onChange(r.entries, id, add)
	}
}

func (r *Registry[T]) Delete(id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.DeleteUnsafe(id)
}

func (r *Registry[T]) DeleteUnsafe(id string) {
	delete(r.entries, id)
	if r.onChange != nil {
		r.onChange(r.entries, id, false)
	}
}

func (r *Registry[T]) Get(id string) (t T, ok bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	t, ok = r.entries[id]
	return
}
