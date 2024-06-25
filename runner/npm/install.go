package npm

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/target"
)

func Install(m target.NodeModule) (err error) {
	if err = exec.Run(m.Abs(), "npm", "install"); err != nil {
		return fmt.Errorf("cannot install node_modules in %s: %s", m.Abs(), err)
	}
	return
}
