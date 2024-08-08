package find

import (
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	. "github.com/cryptopunkscc/portal/target"
)

func Create[T Portal_](
	resolve Resolve[T],
	targets *Cache[T],
	priority Priority,
) Find[T] {
	return FindByPath(source.File, resolve).
		ById(featApps.Path).
		Cached(targets).
		Reduced(priority...)
}
