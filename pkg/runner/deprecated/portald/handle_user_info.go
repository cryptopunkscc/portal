package portald

import (
	user2 "github.com/cryptopunkscc/astrald/mod/user"
)

func (s *Service) UserInfo() (*user2.Info, error) {
	return s.Apphost.User().Info(nil)
}
