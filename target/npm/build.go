package npm

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
)

func Build(_ context.Context, m target.NodeModule) (err error) {
	if err = deps.RequireBinary("npm"); err != nil {
		return
	}
	if !m.PkgJson().CanBuild() {
		return errors.New("missing npm build in package.json")
	}
	if err = exec.Run(m.Abs(), "npm", "run", "build"); err != nil {
		return fmt.Errorf("npm.Build %v: %w", m.Abs(), err)
	}
	return
}
