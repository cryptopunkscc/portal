package npm

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/target"
)

func RunBuild(m target.NodeModule) (err error) {
	if !m.PkgJson().CanBuild() {
		return errors.New("missing npm build in package.json")
	}
	if err = exec.Run(m.Abs(), "npm", "run", "build"); err != nil {
		return fmt.Errorf("npm.RunBuild %v: %w", m.Abs(), err)
	}
	return
}
