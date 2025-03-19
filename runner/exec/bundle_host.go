package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/dir"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"slices"
	"strings"
)

func BundleHostRunner(schemaPrefix ...string) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		defer plog.TraceErr(&err)
		log := plog.Get(ctx).Scope("exec.BundleHostRunner")
		manifest := src.Manifest()
		schemaArr := schemaPrefix
		if manifest.Schema != "" {
			schemaArr = append(schemaArr, manifest.Schema)
		}
		schema := strings.Join(schemaArr, ".")
		log.Println("running:", schema, manifest.Package, args)

		t, err := token.Repository{}.Get(src.Manifest().Package)
		if err != nil {
			return
		}
		args = slices.Insert(args, 0, src.Abs())

		runner, err := target.
			FindByPath(source.File, exec.ResolveBundle).
			OrById(path.Resolver(exec.ResolveBundle, dir.AppSource)).
			Call(ctx, schema)

		if err != nil {
			return
		}
		if len(runner) == 0 {
			return target.ErrNotFound
		}
		runnerBundle := runner[0]

		execFile, err := unpackExecutable(runnerBundle)
		if err != nil {
			return
		}

		return Cmd{}.Run(ctx, t.Token.String(), execFile.Name(), args...)
	}
}
