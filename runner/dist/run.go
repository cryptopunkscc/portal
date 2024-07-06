package dist

import (
	"context"
	"github.com/cryptopunkscc/portal/target"
)

func NewRun(dependencies []target.NodeModule) target.Run[target.Project] {
	return Runner{
		NpmRunner: NewNpmRunner(dependencies),
		GoRunner:  NewGoRunner(),
	}.Run
}

type Runner struct {
	NpmRunner
	GoRunner
}

func (r Runner) Run(ctx context.Context, project target.Project) (err error) {
	switch v := project.(type) {
	case target.ProjectNpm:
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
