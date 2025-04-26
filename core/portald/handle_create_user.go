package portald

import "C"
import "github.com/cryptopunkscc/portal/api/user"

func (s *Service[T]) CreateUser(alias string) (err error) {
	c := user.Client{Client: &s.Apphost}
	if s.UserInfo, err = c.Create(alias); err != nil {
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
