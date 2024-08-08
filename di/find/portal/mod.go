package find

import (
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	. "github.com/cryptopunkscc/portal/target"
)

func Create[T Portal_]() Find[T] {
	return FindByPath(
		source.File,
		sources.Resolver[T]()).
		ById(featApps.Path)
}
