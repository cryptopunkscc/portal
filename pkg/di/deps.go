package di

import (
	"github.com/cryptopunkscc/astrald/sig"
	"reflect"
)

type Create[D any, T any] func(D) T
type Cache struct{ _cache sig.Map[any, any] }
type cache interface{ cache() *sig.Map[any, any] }

func (d *Cache) cache() *sig.Map[any, any] { return &d._cache }
func Single[D any, T any](create Create[D, T], deps any) (t T) {
	key := reflect.ValueOf(create)
	c := deps.(cache).cache()
	if d, ok := c.Get(key); ok {
		return d.(T)
	}
	t = create(deps.(D))
	c.Set(key, t)
	return
}
