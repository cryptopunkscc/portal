package portald

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/astrald"
	"github.com/cryptopunkscc/portal/api/portal"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/api/user"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/resources"
	"github.com/cryptopunkscc/portal/target/app"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/source"
	"path/filepath"
	"sync"
)

type Service[T Portal_] struct {
	cache     Cache[T]
	waitGroup sync.WaitGroup
	processes sig.Map[string, T]
	shutdown  context.CancelFunc

	Config      portal.Config
	configured  bool
	ExtraTokens []string
	AppSources  []Source

	Apphost   apphost.Adapter
	Astrald   astrald.Runner
	Resources resources.FileResources

	Resolve Resolve[Runnable]

	Order []int

	UserCreated *user.Created
}

func (s *Service[T]) Configure() (err error) {
	if err = s.Config.Build(); err != nil {
		return
	}
	if s.Resources, err = resources.NewFileResources(s.Config.Portald, true); err != nil {
		return
	}
	_ = s.ReadCreatedUser()
	s.configured = true
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

func (s *Service[T]) SetupToken(app App_) (err error) {
	_, err = s.Tokens().Resolve(app.Manifest().Package)
	return
}

func (s *Service[T]) Installer() app.Installer {
	return app.Installer{
		Dir:     s.Config.Apps,
		Prepare: s.SetupToken,
	}
}

func (s *Service[T]) Publisher() bundle.Publisher {
	return bundle.Publisher{
		Dir: filepath.Join(s.Config.Astrald, "data"),
	}
}

func (s *Service[T]) Tokens() *token.Repository {
	return token.NewRepository(s.Config.Tokens, &s.Apphost)
}

func (s *Service[T]) apps() Source {
	return source.Dir(s.Config.Apps)
}
