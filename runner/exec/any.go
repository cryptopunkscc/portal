package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runtime/apps"
	"slices"
	"strings"
)

func AnyRunner(cacheDir string, schemaPrefix ...string) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		log := plog.Get(ctx).Scope("exec.AnyRunner")
		manifest := src.Manifest()
		schemaArr := schemaPrefix
		if manifest.Schema != "" {
			schemaArr = append(schemaArr, manifest.Schema)
		}
		schema := strings.Join(schemaArr, ".")
		log.Println("run:", schema, manifest.Package, args)
		runners, err := target.
			FindByPath(source.File, exec2.ResolveBundle).
			ById(path.Resolver(apps.Source)).
			Call(ctx, schema)

		if err != nil {
			return
		}
		if len(runners) == 0 {
			return target.ErrNotFound
		}
		runner := runners[0]
		args = slices.Insert(args, 0, src.Abs())
		err = BundleRunner(cacheDir).Call(ctx, runner, args...)
		return
	}
}
