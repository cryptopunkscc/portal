package dist

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/target"
	"strings"
)

type GoRunner struct {
}

func NewGoRunner() GoRunner {
	return GoRunner{}
}

func (g GoRunner) Run(ctx context.Context, project target.ProjectGo) (err error) {
	cmd := []string{"go", "build", "-o", "dist/main"}
	if project.Manifest().Build != "" {
		cmd = strings.Split(project.Manifest().Build, " ")
	}
	if err = exec.Run(project.Abs(), cmd...); err != nil {
		return fmt.Errorf("cannot install node_modules in %s: %s", project.Abs(), err)
	}
	project.Manifest().Exec = "main"
	return
}
