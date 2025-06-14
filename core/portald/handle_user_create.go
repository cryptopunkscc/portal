package portald

import "C"
import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/user"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (s *Service[T]) CreateUser(alias string) (err error) {
	defer plog.TraceErr(&err)

	// create user
	c := user.Client{Client: &s.Apphost}
	if s.UserCreated, err = c.Create(alias); err != nil {
		return
	}

	// save results
	if err = s.WriteCreatedUser(); err != nil {
		return
	}

	// authenticate as user
	s.Apphost.AuthToken = s.UserCreated.AccessToken
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
	return
}

func (s *Service[T]) signAppContract(identifier string) (err error) {
	id, err := s.Apphost.Resolve(identifier)
	if err != nil {
		return
	}
	c, err := apphost.TokenClient(&s.Apphost).SignAppContract(id)
	if err != nil {
		return
	}
	plog.Scope(identifier).Printf("app contract singed - %s", c)
	return
}

func (s *Service[T]) authenticate() (err error) {
	uat, err := s.Tokens().Resolve("portald")
	if err != nil {
		return
	}
	s.Apphost.AuthToken = uat.Token.String()
	if err = s.Apphost.Reconnect(); err != nil {
		return
	}
	return
}

func (s *Service[T]) PrintCreatedUser() (info *user.Created, err error) {
	if s.UserCreated == nil {
		return nil, plog.Errorf("user not exists")
	}
	return s.UserCreated, nil
}

func (s *Service[T]) WriteCreatedUser() (err error) {
	return s.Resources.WriteYaml("user_created.yaml", s.UserCreated)
}

func (s *Service[T]) ReadCreatedUser() (err error) {
	return s.Resources.ReadYaml("user_created.yaml", &s.UserCreated)
}
