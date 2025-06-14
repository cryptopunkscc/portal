package portald

import "github.com/cryptopunkscc/portal/api/user"

func (s *Service[T]) Claim(alias string) (err error) {
	a := s.Apphost.Clone()
	a.AuthToken = s.UserCreated.AccessToken
	if err = a.Reconnect(); err != nil {
		return
	}

	return user.Client{Client: a}.Claim(alias)
}
