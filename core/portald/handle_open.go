package portald

import (
	"context"
	"errors"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/find"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/supervisor"
)

func (s *Service[T]) Open() Run[apphost.PortaldOpenOpt] {
	return func(ctx context.Context, opt apphost.PortaldOpenOpt, cmd ...string) (err error) {
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
			opt.Order = s.Order
		}

		return find.Runner[T](
			FindByPath(source.File, s.Resolve).
				OrById(path.Resolver(s.Resolve, s.apps())).
				Cached(&s.cache).
				Reduced(Priority{
					Match[Bundle_],
					Match[Dist_],
					Match[Project_]}.
					Sort(opt.Order)...),

			supervisor.Runner[T](
				&s.waitGroup,
				&s.processes,
				multi.Runner[T](s.Runners(schemaPrefix)...),
			),
		).Call(ctx, src, args...)
	}
}
