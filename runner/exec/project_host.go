package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
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

		args = slices.Insert(args, 0, src.Abs())

		runners, err := target.
			FindByPath(source.File, exec2.ResolveProject).
			OrById(path.Resolver(exec2.ResolveProject, env.PortaldApps.Source())).
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
		c, err := exec.Cmd{}.Parse(e, runner.Abs(), strings.Join(args, " "))
		if err != nil {
			return
		}
		log.Println("running", c)
		return Cmd{}.RunApp(ctx, *src.Manifest(), c.Cmd, c.Args...)
	}
}
