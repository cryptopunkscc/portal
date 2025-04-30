package npm_build

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	js "github.com/cryptopunkscc/portal/core/js/embed"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/runner/pack"
	npm2 "github.com/cryptopunkscc/portal/target/npm"
	"slices"
)

func Runner(dependencies ...target.NodeModule) *target.SourceRunner[target.ProjectNpm_] {
	return &target.SourceRunner[target.ProjectNpm_]{
		Resolve: target.Any[target.ProjectNpm_](target.Try(npm2.Resolve_)),
		Runner:  &runner{dependencies: dependencies},
	}
}

type runner struct {
	dependencies []target.NodeModule
}

func NewRun(dependencies ...target.NodeModule) target.Run[target.ProjectNpm_] {
	return (&runner{dependencies: dependencies}).Run
}

func (r *runner) Run(ctx context.Context, project target.ProjectNpm_, args ...string) (err error) {
	plog.Get(ctx).Type(r).Set(&ctx)

	if err = r.setup(); err != nil {
		return
	}
	if r.skip(project, args...) {
		return
	}
	if err = r.prepare(ctx, project); err != nil {
		return
	}
	if err = r.build(ctx, project); err != nil {
		return
	}
	if err = dist.Dist(ctx, project); err != nil {
		return
	}

	if slices.Contains(args, "pack") {
		if err = pack.Run(ctx, project.Dist_()); err != nil {
			return
		}
	}
	return
}

func (r *runner) skip(project target.Project_, args ...string) bool {
	return !project.Changed() && !slices.Contains(args, "clean")
}

func (r *runner) setup() (err error) {
	if r.dependencies == nil {
		r.dependencies = js.LibsDefault
	}
	if len(r.dependencies) == 0 {
		return plog.Errorf("missing js dependencies")
	}
	return
}

func (r *runner) prepare(ctx context.Context, project target.ProjectNpm_) (err error) {
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

func (r *runner) build(ctx context.Context, project target.ProjectNpm_) (err error) {
	plog.Get(ctx).Println("npm run build...")
	if err = npm.Build(ctx, project); err != nil {
		return
	}
	return
}
