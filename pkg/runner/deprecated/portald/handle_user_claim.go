package portald

import (
	"strings"
)

func (s *Service) Claim(alias string) (err error) {
	a := &s.Apphost

	if s.UserCreated != nil {
		a.Token = s.UserCreated.AccessToken.String()
		a.Init()
	}

	id, err := s.Apphost.Resolve(alias)
	if err != nil {
		if len(alias) < 66 && !strings.HasPrefix(alias, ".") {
			var err2 error
			if id, err2 = s.Apphost.Resolve("." + alias); err2 != nil {
				return
			}
		} else {
			return
		}
	}
	sid := id.String()

	_, err = a.User().Claim(nil, sid)
	if err != nil {
		return
	}

	pid, err := a.Resolve("portald")
	if err != nil {
		return
	}

	_, err = a.SignAppContract(pid)
	if err != nil {
		return
	}

	return
}
