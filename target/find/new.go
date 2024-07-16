package find

import (
	embedApps "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/pkg/assets"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/portal"
	"github.com/cryptopunkscc/portal/target/source"
)

type Deps[T target.Portal] interface {
	TargetFinder() target.Finder[T]
	TargetCache() *target.Cache[T]
	Path() target.Path
}

func New[T target.Portal](deps Deps[T]) target.Find[T] {
	launcherSvelteFs := embedApps.LauncherSvelteFS
	resolveEmbedApp := portal.NewResolver[target.App](
		apps.Resolve[target.App](),
		source.FromFS(launcherSvelteFs),
	)
	findPath := target.Mapper[string, string](
		resolveEmbedApp.Path,
		deps.Path(),
	)
	finder := deps.TargetFinder()
	cache := deps.TargetCache()
	return finder.Cached(cache)(
		findPath,
		launcherSvelteFs,
		assets.OsFS{},
	)
}
