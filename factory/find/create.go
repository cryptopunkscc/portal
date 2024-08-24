package find

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
)

func Create[T Portal_](
	path Path,
	targets *Cache[T],
	resolve Resolve[T],
	priority Priority,
) Find[T] {
	return FindByPath(source.File, resolve).
		ById(path).
		Cached(targets).
		Reduced(priority...)
}
