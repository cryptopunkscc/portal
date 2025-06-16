package portald

import "github.com/cryptopunkscc/portal/api/user"

func (s *Service[T]) UserInfo() (*user.Info, error) {
	return user.Op(&s.Apphost).Info()
}
