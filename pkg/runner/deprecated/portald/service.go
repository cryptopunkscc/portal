package portald

import (
	"context"
	"strings"
	"sync"

	"github.com/cryptopunkscc/astrald/mod/user"
	"github.com/cryptopunkscc/portal/pkg/apphost"
	"github.com/cryptopunkscc/portal/pkg/runner/astrald"
	"github.com/cryptopunkscc/portal/pkg/runner/deprecated/portal"
	source2 "github.com/cryptopunkscc/portal/pkg/source"
	app2 "github.com/cryptopunkscc/portal/pkg/source/app"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
	"github.com/cryptopunkscc/portal/pkg/util/resources"
	"github.com/cryptopunkscc/portal/target/app"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/exec"
	"github.com/cryptopunkscc/portal/target/source"
)

type Service struct {
	cache     Cache[Portal_]
	waitGroup sync.WaitGroup
	shutdown  context.CancelFunc

	Config      portal.Config
	configured  bool
	ExtraTokens []string
	AppSources  []Source

	Apphost   client.Adapter
	Astrald   astrald.Runner
	Resources resources.Dir

	Resolve Resolve[Runnable]

	Order []int

	UserCreated *user.CreatedUserInfo
	User        *user.Info
	hasUser     bool
	*NodeInfo
}

func (s *Service) Configure() (err error) {
	if err = s.Config.Build(); err != nil {
		return
	}
	s.Resources.Path = s.Config.Portald
	if err = s.Resources.Init(); err != nil {
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

func (s *Service) user() *user.Info {
	if s.User == nil {
		s.User, _ = s.UserInfo()
	}
	return s.User
}

func (s *Service) HasUser() bool {
	if !s.hasUser {
		_, err := s.UserInfo()
		s.hasUser = err == nil || strings.Contains(err.Error(), "(1)")
	}
	return s.hasUser
}

func (s *Service) Installer() app.Installer {
	return app.Installer{
		Dir: s.Config.Apps,
		Repositories: Repositories{
			source.Repository,
			s.Bundles(),
		},
		Resolvers: []Resolve[Source]{
			exec.ResolveDist.Try,
			exec.ResolveBundle.Try,
		},
		Prepare: s.ClaimApp,
	}
}

func (s *Service) Publisher() app2.Publisher {
	return app2.Publisher{ClientObjects: &s.Apphost.Objects().ObjectsClient}
}

func (s *Service) Tokens() *client.Tokens {
	return s.Apphost.Tokens(s.Config.Tokens)
}

func (s *Service) apps() Source {
	return source.Dir(s.Config.Apps)
}

func (s *Service) appsRef() *source2.Ref {
	return source2.OSRef(s.Config.Apps)
}

func (s *Service) Bundles() *bundle.Repository {
	return &bundle.Repository{Apphost: &s.Apphost}
}

func (s *Service) AppObjects() *app2.Objects {
	return &app2.Objects{Adapter: &s.Apphost}
}
