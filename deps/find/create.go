package find

import (
	embedApps "github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/mock/appstore"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/portal"
	"github.com/cryptopunkscc/go-astral-js/target/source"
)

func Create[T target.Portal](
	cache *target.Cache[T],
	finder target.Finder[T],
) target.Find[T] {
	resolveEmbed := portal.NewResolver[target.App](
		apps.Resolve[target.App](),
		source.FromFS(embedApps.LauncherSvelteFS),
	)
	findPath := target.Mapper[string, string](
		resolveEmbed.Path,
		appstore.Path,
	)

	return finder.Cached(cache)(findPath, embedApps.LauncherSvelteFS)
}
