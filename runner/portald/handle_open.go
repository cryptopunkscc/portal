package portald

import (
	"context"
	"errors"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/unknown"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/find"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/supervisor"
	"github.com/cryptopunkscc/portal/runtime/apps"
)

func (s *Runner[T]) Open() Run[portald.OpenOpt] {
	return func(ctx context.Context, opt portald.OpenOpt, cmd ...string) (err error) {
		if len(cmd) == 0 {
			return errors.New("no command")
		}
		src := cmd[0]
		args := cmd[1:]

		var schemaPrefix []string
		if opt.Schema != "" {
			schemaPrefix = []string{opt.Schema}
		}
		plog.Get(ctx).Type(s).Println("open:", opt, cmd, opt.Order)

		order := Priority{
			Match[Bundle_],
			Match[Dist_],
			Match[Project_],
		}.Sort(opt.Order)

		return find.Runner[T](
			FindByPath(
				source.File, Any[T](
					Skip("node_modules"),
					Try(exec2.ResolveBundle),
					Try(exec2.ResolveDist),
					Try(unknown.ResolveBundle),
					Try(unknown.ResolveDist),
					Try(unknown.ResolveProject),
				),
			).ById(path.Resolver(apps.Source)).
				Cached(&s.cache).
				Reduced(order...),
			supervisor.Runner[T](
				&s.waitGroup,
				&s.processes,
				multi.Runner[T](
					app.Runner(exec.BundleRunner(s.CacheDir)),
					app.Runner(exec.DistRunner()),
					app.Runner(exec.AnyRunner(s.CacheDir, schemaPrefix...)),
				),
			),
		).Call(ctx, src, args...)
	}
}
