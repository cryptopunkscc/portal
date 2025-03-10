package portald

import (
	"context"
	"errors"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/portal"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/unknown"
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

		if len(opt.Order) == 0 {
			opt.Order = defaultOrder
		}

		order := Priority{
			Match[Bundle_],
			Match[Dist_],
			Match[Project_],
		}.Sort(opt.Order)

		return find.Runner[T](
			FindByPath(
				source.File, Any[T](
					Skip("node_modules"),
					Try(exec.ResolveBundle),
					Try(exec.ResolveDist),
					Try(exec.ResolveProject),
					Try(unknown.ResolveBundle),
					Try(unknown.ResolveDist),
					Try(unknown.ResolveProject),
				),
			).ById(path.Resolver(Any[T](
				Try(portal.Resolve_),
				Try(exec.ResolveBundle),
			), apps.Source)).
				Cached(&s.cache).
				Reduced(order...),
			supervisor.Runner[T](
				&s.waitGroup,
				&s.processes,
				multi.Runner[T](s.runners(schemaPrefix)...),
			),
		).Call(ctx, src, args...)
	}
}
