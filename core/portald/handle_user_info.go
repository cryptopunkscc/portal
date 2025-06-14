package portald

import "github.com/cryptopunkscc/portal/api/user"

func (s *Service[T]) UserInfo() (*user.Info, error) {
	return user.Client{Rpc: s.Apphost.Rpc()}.Info()
}
