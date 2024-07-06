package dist

import (
	"context"
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/target"
)

type NpmRunner struct {
	dependencies []target.NodeModule
}

func NewNpmRunner(dependencies []target.NodeModule) NpmRunner {
	return NpmRunner{dependencies: dependencies}
}

func (r NpmRunner) Run(ctx context.Context, project target.ProjectNpm) (err error) {
	if err = r.Prepare(ctx, project); err != nil {
		return
	}
	if err = r.Build(ctx, project); err != nil {
		return
	}
	if err = Dist(ctx, project); err != nil {
		return
	}
	return
}

func (r NpmRunner) Prepare(ctx context.Context, project target.ProjectNpm) (err error) {
	if err = npm.Install(project); err != nil {
		return
	}
	if err = npm.NewInjector(r.dependencies).Run(ctx, project); err != nil {
		return
	}
	return
}

func (r NpmRunner) Build(ctx context.Context, project target.ProjectNpm) (err error) {
	if err = npm.RunBuild(project); err != nil {
		return
	}
	return
}
