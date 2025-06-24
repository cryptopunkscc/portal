package portald

import "github.com/cryptopunkscc/portal/api/user"

func (s *Service) UserInfo() (*user.Info, error) {
	return user.Op(&s.Apphost).Info()
}
