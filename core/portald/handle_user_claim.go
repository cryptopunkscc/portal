package portald

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/dir"
	"github.com/cryptopunkscc/portal/api/user"
)

func (s *Service) Claim(alias string) (err error) {
	a := s.Apphost.Clone()
	a.AuthToken = s.UserCreated.AccessToken
	if err = a.Reconnect(); err != nil {
		return
	}

	id, err := s.Apphost.Resolve(alias)
	if err != nil {
		return
	}
	sid := id.String()

	err = user.Op(a).Claim(sid)
	if err != nil {
		return
	}

	pid, err := dir.Op(a, sid).Resolve("portald")
	if err != nil {
		return
	}

	_, err = apphost.Op(a, sid).SignAppContract(pid)
	if err != nil {
		return
	}

	return
}
