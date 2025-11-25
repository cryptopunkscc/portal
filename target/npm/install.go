package npm

import (
	"context"
	"fmt"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
)

func Install(_ context.Context, m target.NodeModule) (err error) {
	if err = deps.Check("npm", "-v"); err != nil {
		return
	}
	if err = exec.Run(m.Abs(), "npm", "install"); err != nil {
		return fmt.Errorf("cannot install node_modules in %s: %s", m.Abs(), err)
	}
	return
}
