package apps

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/require"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
)

func Path(sources ...target.Source) target.Path {
	if len(sources) == 0 {
		sources = []target.Source{
			require.NoErr(source.File(DefaultDir())),
		}
	}
	return func(ctx context.Context, port string) (path string, err error) {
		log := plog.Get(ctx)
		log.Println("resolving path", port)
		for _, t := range apps.Resolver[target.Bundle_]().List(sources...) {
			log.Println("matching target", t.Manifest())
			if t.Manifest().Match(port) {
				path = t.Abs()
				return
			}
		}
		err = target.ErrNotFound
		return
	}
}
