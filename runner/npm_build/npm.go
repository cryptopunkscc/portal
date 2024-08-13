package npm_build

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/target"
)

type runner struct {
	dependencies []target.NodeModule
}

func Runner(dependencies ...target.NodeModule) target.Run[target.ProjectNpm_] {
	return runner{dependencies: dependencies}.Run
}

func (r runner) Run(ctx context.Context, project target.ProjectNpm_) (err error) {
	plog.Get(ctx).Type(r).Set(&ctx)
	if err = r.prepare(ctx, project); err != nil {
		return
	}
	if err = r.build(ctx, project); err != nil {
		return
	}
	if err = dist.Dist(ctx, project); err != nil {
		return
	}
	return
}

func (r runner) prepare(ctx context.Context, project target.ProjectNpm_) (err error) {
	log := plog.Get(ctx)
	log.Println("npm install...")
	if err = npm.Install(ctx, project); err != nil {
		return
	}
	for i, dependency := range r.dependencies {
		log.Println(i, dependency.Abs(), dependency.Path())
	}
	log.Println("injecting portal lib...")
	inject := npm.Injector(r.dependencies)
	if err = inject(ctx, project); err != nil {
		return
	}
	return
}

func (r runner) build(ctx context.Context, project target.ProjectNpm_) (err error) {
	plog.Get(ctx).Println("npm run build...")
	if err = npm.Build(ctx, project); err != nil {
		return
	}
	return
}
