package find

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
)

func Create[T Portal_]() Find[T] {
	return FindByPath(
		source.File,
		sources.Resolver[T]()).
		ById(appstore.Path)
}
