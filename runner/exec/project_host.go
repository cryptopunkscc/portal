package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runtime/dir"
	"github.com/cryptopunkscc/portal/runtime/tokens"
	"slices"
	"strings"
)

func ProjectHostRunner(schemaPrefix ...string) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		defer plog.TraceErr(&err)
		log := plog.Get(ctx).Scope("exec.ProjectHostRunner")
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

		runners, err := target.
			FindByPath(source.File, exec2.ResolveProject).
			OrById(path.Resolver(exec2.ResolveProject, dir.AppSource)).
			Call(ctx, schema)

		if err != nil {
			return
		}
		if len(runners) == 0 {
			return target.ErrNotFound
		}
		var runner target.ProjectExec
		for _, r := range runners {
			if r.Manifest().Exec != "" {
				runner = r
				break
			}
		}
		if runner == nil {
			return target.ErrNotFound
		}

		e := runner.Manifest().Exec
		c := exec.Cmd{}.Parse(e, runner.Abs(), strings.Join(args, " "))
		log.Println("running", c)
		return Cmd{}.Run(ctx, token.Token.String(), c.Cmd, c.Args...)
	}
}
