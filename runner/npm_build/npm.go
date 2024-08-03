package npm_build

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/target"
)

type NpmRunner struct {
	dependencies []target.NodeModule
}

func NewRunner(dependencies []target.NodeModule) NpmRunner {
	return NpmRunner{dependencies: dependencies}
}

func (r NpmRunner) Run(ctx context.Context, project target.ProjectNpm_) (err error) {
	plog.Get(ctx).Type(r).Set(&ctx)
	if err = r.Prepare(ctx, project); err != nil {
		return
	}
	if err = r.Build(ctx, project); err != nil {
		return
	}
	if err = dist.Dist(ctx, project); err != nil {
		return
	}
	return
}

func (r NpmRunner) Prepare(ctx context.Context, project target.ProjectNpm_) (err error) {
	log := plog.Get(ctx)
	log.Println("npm install...")
	if err = npm.Install(project); err != nil {
		return
	}
	for i, dependency := range r.dependencies {
		log.Println(i, dependency.Abs(), dependency.Path())
	}
	log.Println("injecting portal lib...")
	if err = npm.NewInjector(r.dependencies).Run(ctx, project); err != nil {
		return
	}
	return
}

func (r NpmRunner) Build(ctx context.Context, project target.ProjectNpm_) (err error) {
	plog.Get(ctx).Println("npm run build...")
	if err = npm.RunBuild(project); err != nil {
		return
	}
	return
}
