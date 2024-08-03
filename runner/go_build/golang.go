package go_build

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/target"
	"strings"
)

type GoRunner struct {
}

func NewRunner() GoRunner {
	return GoRunner{}
}

func (g GoRunner) Run(ctx context.Context, project target.ProjectGo) (err error) {
	if err = deps.RequireBinary("go"); err != nil {
		return
	}
	cmd := []string{"go", "build", "-o", "dist/main"}
	if project.Manifest().Build != "" {
		cmd = strings.Split(project.Manifest().Build, " ")
	}
	if err = exec.Run(project.Abs(), cmd...); err != nil {
		return fmt.Errorf("run golang build %s: %s", project.Abs(), err)
	}
	project.Manifest().Exec = "main"
	if err = dist.Dist(ctx, project); err != nil {
		return
	}
	return
}
