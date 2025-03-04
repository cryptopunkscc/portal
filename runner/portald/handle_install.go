package portald

import (
	"github.com/cryptopunkscc/portal/runner/install"
	apps2 "github.com/cryptopunkscc/portal/runtime/apps"
)

func (s *Runner[T]) Install(src string) (c <-chan install.Result, err error) {
	return install.Runner(apps2.Dir).Run(src)
}
