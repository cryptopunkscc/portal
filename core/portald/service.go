package portald

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/astrald"
	"github.com/cryptopunkscc/portal/api/portal"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/install"
	"sync"
)

type Service[T Portal_] struct {
	cache     Cache[T]
	waitGroup sync.WaitGroup
	processes sig.Map[string, T]
	shutdown  context.CancelFunc

	Config      portal.Config
	ExtraTokens []string

	Apphost apphost.Adapter
	Astrald astrald.Runner

	Resolve Resolve[T]
	Runners func([]string) []Run[Portal_]
	Order   []int
}

func (s *Service[T]) Configure() (err error) {
	err = s.Config.Build()
	plog.D().Printf("config:\n%s", s.Config.Yaml())
	return
}

func (s *Service[T]) Stop() {
	s.shutdown()
}

func (s *Service[T]) Wait() (err error) {
	s.waitGroup.Wait()
	s.shutdown()
	return
}

func (s *Service[T]) Install() install.Runner {
	return install.Runner{
		AppsDir: s.Config.Apps,
		Tokens:  *s.Tokens(),
	}
}

func (s *Service[T]) Tokens() *token.Repository {
	return token.NewRepository(s.Config.Tokens, &s.Apphost)
}

func (s *Service[T]) apps() Source {
	return source.Dir(s.Config.Apps)
}
