package find

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/resolve/source"
)

func Create[T Portal_](
	targets *Cache[T],
	resolve Resolve[T],
	priority Priority,
) Find[T] {
	return FindByPath(source.File, resolve).
		ById(appstore.Path).
		Cached(targets).
		Reduced(priority...)
}
