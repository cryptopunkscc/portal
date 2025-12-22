package portald

import (
	user2 "github.com/cryptopunkscc/astrald/mod/user"
	"github.com/cryptopunkscc/portal/api/user"
)

func (s *Service) UserInfo() (*user2.Info, error) {
	return user.Op(&s.Apphost).Info()
}
