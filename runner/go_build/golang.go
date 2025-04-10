package go_build

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	golang "github.com/cryptopunkscc/portal/resolve/go"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/pack"
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

func Runner(platforms ...string) *target.SourceRunner[target.ProjectGo] {
	return &target.SourceRunner[target.ProjectGo]{
		Resolve: target.Any[target.ProjectGo](target.Try(golang.ResolveProject)),
		Runner:  runner{platforms: platforms},
	}
}

type runner struct{ platforms []string }

func NewRun(platforms ...string) target.Run[target.ProjectGo] {
	return runner{platforms}.Run
}

func (g runner) Run(ctx context.Context, project target.ProjectGo, args ...string) (err error) {
	log := plog.Get(ctx).Type(g).Set(&ctx)
	if err = deps.RequireBinary("go"); err != nil {
		return
	}

	if !project.Changed() && !slices.Contains(args, "clean") {
		return
	}

	if len(g.platforms) == 0 {
		g.platforms = []string{runtime.GOOS}
	}

	log.Printf("go build %T %s %v", project, project.Abs(), g.platforms)
	cmd := exec.Cmd{
		Cmd:  "go",
		Args: []string{"build", "-o", "dist/main"},
		Dir:  project.Abs(),
	}.Default()

	if slices.Contains(args, "clean") {
		if err = os.RemoveAll(filepath.Join(project.Abs(), "dist")); err != nil {
			log.W().Println(err)
		}
	}
	for _, platform := range g.platforms {
		build, ok := project.Build()[platform]
		if !ok {
			build, ok = project.Build()["default"]
		}
		if ok {
			if cmd, err = cmd.Parse(build.Cmd); err != nil {
				return
			}
			cmd = cmd.AddEnv(build.Env...).AddEnv("GOOS=" + platform)
		}
		if err = cmd.Build().Run(); err != nil {
			return fmt.Errorf("run golang build %s: %s", project.Abs(), err)
		}
		project.Manifest().Exec = build.Out
		if err = dist.Dist(ctx, project); err != nil {
			return
		}

		if slices.Contains(args, "pack") {
			if err = pack.Run(ctx, project.Dist_()); err != nil {
				return
			}
		}
	}
	return
}
