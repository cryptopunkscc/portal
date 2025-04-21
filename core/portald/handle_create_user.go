package portald

func (s *Service[T]) CreateUser(alias string) (err error) {
	if s.UserInfo, err = s.Apphost.User().Create(alias); err != nil {
		return
	}
	return s.WriteUserInfo()
}

func (s *Service[T]) WriteUserInfo() (err error) {
	return s.Resources.WriteYaml("user_info.yaml", s.UserInfo)
}

func (s *Service[T]) ReadUserInfo() (err error) {
	return s.Resources.ReadYaml("user_info.yaml", &s.UserInfo)
}
