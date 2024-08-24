package find

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/require"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runtime/apps"
)

func Create[T Portal_]() Find[T] {
	return FindByPath(
		source.File,
		sources.Resolver[T]()).
		ById(apps.Path(require.NoErr(source.File(apps.DefaultDir()))))
}
