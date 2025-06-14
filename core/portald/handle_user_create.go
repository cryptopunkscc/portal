package portald

import "C"
import (
	"github.com/cryptopunkscc/portal/api/user"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (s *Service[T]) CreateUser(alias string) (err error) {
	c := user.Client{Client: &s.Apphost}
	if s.UserCreated, err = c.Create(alias); err != nil {
		return
	}
	return s.WriteCreatedUser()
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
