package portald

import (
	user2 "github.com/cryptopunkscc/astrald/mod/user"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (s *Service) CreateUser(alias string) (err error) {
	defer plog.TraceErr(&err)

	// create user
	if s.UserCreated, err = s.Apphost.User().Create(nil, alias); err != nil {
		return
	}

	// save results
	if err = s.WriteCreatedUser(); err != nil {
		return
	}

	// authenticate as user
	s.Apphost.Token = s.UserCreated.AccessToken.String()
	if err = s.Apphost.Reconnect(); err != nil {
		return
	}

	// sign portald
	err = s.signAppContract("portald")
	if err != nil {
		return
	}

	// authenticate as portald
	err = s.authenticate()

	s.hasUser = s.HasUser()

	// sign installed apps
	for _, app := range s.InstalledApps(ListAppsOpts{Hidden: true}) {
		_, _ = s.ClaimPackage(app.Package)
	}
	return
}

func (s *Service) signAppContract(identifier string) (err error) {
	id, err := s.Apphost.Resolve(identifier)
	if err != nil {
		return
	}
	c, err := s.Apphost.SignAppContract(id)
	if err != nil {
		return
	}
	plog.Scope(identifier).Printf("app contract singed - %s", c)
	return
}

func (s *Service) authenticate() (err error) {
	uat, err := s.Tokens().Resolve("portald")
	if err != nil {
		return
	}
	s.Apphost.Token = uat.Token.String()
	if err = s.Apphost.Reconnect(); err != nil {
		return
	}
	return
}

func (s *Service) PrintCreatedUser() (info *user2.CreatedUserInfo, err error) {
	if s.UserCreated == nil {
		return nil, plog.Errorf("user not exists")
	}
	return s.UserCreated, nil
}

func (s *Service) WriteCreatedUser() (err error) {
	return s.Resources.WriteYaml("user_created.yaml", s.UserCreated)
}

func (s *Service) ReadCreatedUser() (err error) {
	return s.Resources.ReadYaml("user_created.yaml", &s.UserCreated)
}
