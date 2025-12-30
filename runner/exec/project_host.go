package exec

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/cryptopunkscc/portal/api/portald"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/dist"
	exec2 "github.com/cryptopunkscc/portal/target/exec"
	"github.com/cryptopunkscc/portal/target/project"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/google/shlex"
)

func (r Runner) ProjectHost() *target.SourceRunner[target.Portal_] {
	return &target.SourceRunner[target.Portal_]{
		Resolve: target.Any[target.Portal_](
			target.Try(dist.Resolve_),
			target.Try(project.Resolve_),
		),
		Runner: &ProjectHostRunner{
			Runner: r,
			Repo: target.SourcesRepository[target.Project_]{
				Sources: []target.Source{source.Dir(r.Apps)},
				Resolve: target.Any[target.Project_](exec2.ResolveProject.Try),
			},
		},
	}
}

type ProjectHostRunner struct {
	Runner
	Repo target.SourcesRepository[target.Project_]
}

func (r *ProjectHostRunner) Run(ctx context.Context, src target.Portal_, args ...string) (err error) {
	defer plog.TraceErr(&err)
	if src.Manifest().Runtime == "" {
		return errors.New("ProjectHostRunner requires a schema declared in manifest")
	}

	log := plog.Get(ctx).Type(r)

	runtime := src.Manifest().Runtime
	opt := portald.OpenOpt{}
	if opt.Load(ctx); len(opt.Schema) > 0 {
		runtime = opt.Schema + "." + runtime
	}

	log.Println("running:", runtime, src.Manifest(), args)

	hostBundle := r.Repo.First(runtime)
	if hostBundle == nil {
		return fmt.Errorf("could not runner for %s", runtime)
	}

	b := hostBundle.Build().Get()
	if b.Exec == "" {
		return target.ErrNotFound
	}

	cmd, err := shlex.Split(b.Exec)
	if err != nil {
		return
	}

	args = slices.Concat(cmd[1:], []string{hostBundle.Abs(), src.Abs()}, args)

	return r.RunApp(ctx, *src.Manifest(), cmd[0], args...)
}
