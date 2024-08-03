package all_build

import (
	"context"
	"github.com/cryptopunkscc/portal/runner/go_build"
	"github.com/cryptopunkscc/portal/runner/npm_build"
	"github.com/cryptopunkscc/portal/target"
)

func NewRun(dependencies []target.NodeModule) target.Run[target.Project_] {
	return Runner{
		NpmRunner: npm_build.NewRunner(dependencies),
		GoRunner:  go_build.NewRunner(),
	}.Run
}

type Runner struct {
	npm_build.NpmRunner
	go_build.GoRunner
}

func (r Runner) Run(ctx context.Context, project target.Project_) (err error) {
	switch v := project.(type) {
	case target.ProjectNpm_:
		if err = r.NpmRunner.Run(ctx, v); err != nil {
			return
		}
	case target.ProjectGo:
		if err = r.GoRunner.Run(ctx, v); err != nil {
			return
		}
	}
	return
}
