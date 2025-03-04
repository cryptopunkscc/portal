package portald

import (
	"github.com/cryptopunkscc/portal/runner/uninstall"
	apps2 "github.com/cryptopunkscc/portal/runtime/apps"
)

func (s *Runner[T]) Uninstall(app string) error {
	return uninstall.Runner(apps2.Source)(app)
}
