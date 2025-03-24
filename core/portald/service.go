package portald

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/astrald"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/install"
	"sync"
)

type Service[T Portal_] struct {
	cache     Cache[T]
	waitGroup sync.WaitGroup
	processes sig.Map[string, T]
	shutdown  context.CancelFunc

	CreateTokens []string
	NodeDir      mem.String
	AppsDir      mem.String
	TokensDir    mem.String

	Apphost apphost.Adapter
	Astrald astrald.Runner

	Resolve Resolve[T]
	Runners func([]string) []Run[Portal_]
	Order   []int
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
		AppsDir: s.AppsDir,
		Tokens:  *s.Tokens(),
	}
}

func (s *Service[T]) Tokens() *token.Repository {
	return token.NewRepository(s.TokensDir, &s.Apphost)
}

func (s *Service[T]) apps() Source {
	return source.Dir(s.AppsDir.Require())
}
