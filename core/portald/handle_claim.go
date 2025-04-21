package portald

func (s *Service[T]) Claim(alias string) (err error) {
	a := s.Apphost.Clone()
	a.AuthToken = s.UserInfo.AccessToken
	if err = a.Reconnect(); err != nil {
		return
	}
	return a.User().Claim(alias)
}
