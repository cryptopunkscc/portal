package npm

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/project"
	"slices"
)

func BuildRunner(dependencies ...target.NodeModule) *target.SourceRunner[target.ProjectNpm_] {
	return &target.SourceRunner[target.ProjectNpm_]{
		Resolve: target.Any[target.ProjectNpm_](target.Try(Resolve_)),
		Runner:  &buildRunner{dependencies: dependencies},
	}
}

func BuildProject(dependencies ...target.NodeModule) target.Run[target.ProjectNpm_] {
	return (&buildRunner{dependencies: dependencies}).Run
}

type buildRunner struct {
	dependencies []target.NodeModule
}

func (r *buildRunner) Run(ctx context.Context, projectNpm target.ProjectNpm_, args ...string) (err error) {
	plog.Get(ctx).Type(r).Set(&ctx)

	if err = r.setup(); err != nil {
		return
	}
	if r.skip(projectNpm, args...) {
		return
	}
	if err = r.prepare(ctx, projectNpm); err != nil {
		return
	}
	if err = r.build(ctx, projectNpm); err != nil {
		return
	}
	if err = project.Dist(ctx, projectNpm); err != nil {
		return
	}

	if slices.Contains(args, "pack") {
		if err = dist.Pack(projectNpm.Dist_()); err != nil {
			return
		}
	}
	return
}

func (r *buildRunner) skip(project target.Project_, args ...string) bool {
	return !project.Changed() && !slices.Contains(args, "clean")
}

func (r *buildRunner) setup() (err error) {
	if r.dependencies == nil {
		r.dependencies = LibsDefault
	}
	if len(r.dependencies) == 0 {
		return plog.Errorf("missing js dependencies")
	}
	return
}

func (r *buildRunner) prepare(ctx context.Context, project target.ProjectNpm_) (err error) {
	log := plog.Get(ctx)
	log.Println("npm install...")
	if err = Install(ctx, project); err != nil {
		return
	}
	for i, dependency := range r.dependencies {
		log.Println(i, dependency.Abs(), dependency.Path())
	}
	log.Println("injecting portal lib...")
	inject := Injector(r.dependencies)
	if err = inject(ctx, project); err != nil {
		return
	}
	return
}

func (r *buildRunner) build(ctx context.Context, project target.ProjectNpm_) (err error) {
	plog.Get(ctx).Println("npm run build...")
	if err = BuildModule(ctx, project); err != nil {
		return
	}
	return
}

func BuildModule(_ context.Context, m target.NodeModule) (err error) {
	if err = deps.RequireBinary("npm"); err != nil {
		return
	}
	if !m.PkgJson().CanBuild() {
		return errors.New("missing npm build in package.json")
	}
	if err = exec.Run(m.Abs(), "npm", "run", "build"); err != nil {
		return fmt.Errorf("npm.BuildModule %v: %w", m.Abs(), err)
	}
	return
}
