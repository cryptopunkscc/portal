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

type Service struct {
	cache     Cache[Portal_]
	waitGroup sync.WaitGroup
	processes sig.Map[string, Portal_]
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
	User        *user.Info
	*NodeInfo
}

func (s *Service) Configure() (err error) {
	if err = s.Config.Build(); err != nil {
		return
	}
	if s.Resources, err = resources.NewFileResources(s.Config.Portald, true); err != nil {
		return
	}
	_ = s.ReadCreatedUser()
	s.configured = true
	plog.D().Scope("config").Printf("\n%s", s.Config.Yaml())
	return
}

func (s *Service) Stop() {
	s.shutdown()
}

func (s *Service) Wait() (err error) {
	s.waitGroup.Wait()
	s.shutdown()
	return
}

func (s *Service) SetupToken(app App_) (err error) {
	_, err = s.Tokens().Resolve(app.Manifest().Package)
	return
}

func (s *Service) PrepareApp(app App_) (err error) {
	t, err := s.Tokens().Resolve(app.Manifest().Package)
	if err != nil {
		return
	}
	if s.User != nil {
		err = s.signAppContract(t.Identity.String())
		if err != nil {
			return
		}
	}
	return
}

func (s *Service) Installer() app.Installer {
	return app.Installer{
		Dir:     s.Config.Apps,
		Prepare: s.PrepareApp,
	}
}

func (s *Service) Publisher() bundle.Publisher {
	return bundle.Publisher{
		Dir: filepath.Join(s.Config.Astrald, "data"),
	}
}

func (s *Service) Tokens() *token.Repository {
	return token.NewRepository(s.Config.Tokens, &s.Apphost)
}

func (s *Service) apps() Source {
	return source.Dir(s.Config.Apps)
}
