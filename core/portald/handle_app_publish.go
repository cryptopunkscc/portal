package portald

import (
	"github.com/cryptopunkscc/portal/source/app"
)

func (s *Service) PublishApps(path string) (out []app.ReleaseInfo, err error) {
	return s.Publisher2().PublishBundles(path)
}
