package unknown

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/dist"
	"github.com/cryptopunkscc/portal/resolve/project"
)

var ResolveDist Resolve[Dist_] = func(src Source) (Dist_, error) { return dist.Resolver[any](resolve)(src) }
var ResolveBundle Resolve[Bundle_] = func(src Source) (Bundle_, error) { return bundle.Resolver[any](dist.Resolver[any](resolve))(src) }
var ResolveProject Resolve[Project_] = func(src Source) (Project_, error) { return project.Resolver[any](resolve)(src) }

func resolve(Source) (t any, err error) { return }
