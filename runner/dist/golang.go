package dist

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type GoRunner struct {
}

func NewGoRunner() GoRunner {
	return GoRunner{}
}

func (g GoRunner) Run(ctx context.Context, project target.ProjectGo) (err error) {
	if err = exec.Run(project.Abs(), "go", "build", "-o", "dist/main"); err != nil {
		return fmt.Errorf("cannot install node_modules in %s: %s", project.Abs(), err)
	}
	project.Manifest().Exec = "main"
	return
}
