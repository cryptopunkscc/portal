package di

import (
	"github.com/cryptopunkscc/astrald/sig"
	"reflect"
)

type Create[D any, T any] func(D) T
type Cache struct{ _cache sig.Map[any, any] }
type cache interface{ cache() *sig.Map[any, any] }

func (d *Cache) cache() *sig.Map[any, any]              { return &d._cache }
func S[D any, T any](create Create[D, T], deps D) (t T) { return Singleton[D, T](create, deps) }
func Singleton[D any, T any](create Create[D, T], deps D) (t T) {
	key := reflect.ValueOf(create)
	c := any(deps).(cache).cache()
	if d, ok := c.Get(key); ok {
		return d.(T)
	}
	t = create(deps)
	c.Set(key, t)
	return
}
