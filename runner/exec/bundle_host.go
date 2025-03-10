package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runtime/apps"
	"github.com/cryptopunkscc/portal/runtime/tokens"
	"slices"
	"strings"
)

func BundleHostRunner(cacheDir string, schemaPrefix ...string) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		log := plog.Get(ctx).Scope("exec.BundleHostRunner")
		manifest := src.Manifest()
		schemaArr := schemaPrefix
		if manifest.Schema != "" {
			schemaArr = append(schemaArr, manifest.Schema)
		}
		schema := strings.Join(schemaArr, ".")
		log.Println("running:", schema, manifest.Package, args)

		token, err := tokens.Repository{}.Get(src.Manifest().Package)
		if err != nil {
			return
		}
		args = slices.Insert(args, 0, src.Abs())

		runner, err := target.
			FindByPath(source.File, exec.ResolveBundle).
			ById(path.Resolver(exec.ResolveBundle, apps.Source)).
			Call(ctx, schema)

		if err != nil {
			return
		}
		if len(runner) == 0 {
			return target.ErrNotFound
		}
		runnerBundle := runner[0]

		execFile, err := unpackExecutable(cacheDir, runnerBundle)
		if err != nil {
			return
		}

		return Cmd{}.Run(ctx, token.Token.String(), execFile.Name(), args...)
	}
}
