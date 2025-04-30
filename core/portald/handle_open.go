package portald

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/portald"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/portal"
	"github.com/cryptopunkscc/portal/target/source"
)

func (s *Service[T]) Open() Run[portald.OpenOpt] {
	dispatcher := s.dispatcher()
	return func(ctx context.Context, opt portald.OpenOpt, cmd ...string) (err error) {
		plog.Get(ctx).Type(s).Println("open:", opt, cmd)
		if len(cmd) == 0 {
			return errors.New("no command")
		}
		if len(opt.Order) == 0 {
			opt.Order = s.Order
		}
		p := dispatcher
		p.Priority = p.Sort(opt.Order)
		src := cmd[0]
		args := cmd[1:]
		opt.Save(&ctx)
		return p.Run(ctx, src, args...)
	}
}

func (s *Service[T]) dispatcher() Dispatcher {
	return Dispatcher{
		Provider: s.Provider(),
		Runner: &CachedRunner[T]{
			Cache: &s.cache,
			Runner: &AsyncRunner{
				WaitGroup: &s.waitGroup,
			},
		},
	}
}

func (s *Service[T]) Provider() Provider[Runnable] {
	return Provider[Runnable]{
		Priority: Priority{
			Match[Bundle_],
			Match[Dist_],
			Match[Project_],
		},
		Repository: Repositories{
			s.cache.Repository(),
			portal.Repository(s.apps()),
			source.Repository,
		},
		Resolve: s.Resolve,
	}
}
