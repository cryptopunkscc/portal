package portald

import (
	"github.com/cryptopunkscc/portal/pkg/source/app"
)

func (s *Service) PublishApps(path string) (out []app.ReleaseInfo, err error) {
	return s.Publisher().PublishBundles(path)
}
