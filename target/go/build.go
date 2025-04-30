package golang

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/project"
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

func BuildRunner(platforms ...string) *target.SourceRunner[target.ProjectGo] {
	return &target.SourceRunner[target.ProjectGo]{
		Resolve: ResolveProject,
		Runner:  buildRunner{platforms: platforms},
	}
}

func BuildProject(platforms ...string) target.Run[target.ProjectGo] {
	return buildRunner{platforms}.Run
}

type buildRunner struct{ platforms []string }

func (g buildRunner) Run(ctx context.Context, projectGo target.ProjectGo, args ...string) (err error) {
	log := plog.Get(ctx).Type(g).Set(&ctx)
	if err = deps.RequireBinary("go"); err != nil {
		return
	}

	if !projectGo.Changed() && !slices.Contains(args, "clean") {
		return
	}

	if len(g.platforms) == 0 {
		g.platforms = []string{runtime.GOOS}
	}

	log.Printf("go build %T %s %v", projectGo, projectGo.Abs(), g.platforms)
	cmd := exec.Cmd{
		Cmd:  "go",
		Args: []string{"build", "-o", "dist/main"},
		Dir:  projectGo.Abs(),
	}.Default()

	if slices.Contains(args, "clean") {
		if err = os.RemoveAll(filepath.Join(projectGo.Abs(), "dist")); err != nil {
			log.W().Println(err)
		}
	}
	for _, platform := range g.platforms {
		build, ok := projectGo.Build()[platform]
		if !ok {
			build, ok = projectGo.Build()["default"]
		}
		if ok {
			if cmd, err = cmd.Parse(build.Cmd); err != nil {
				return
			}
			cmd = cmd.AddEnv(build.Env...).AddEnv("GOOS=" + platform)
		}
		if err = cmd.Build().Run(); err != nil {
			return fmt.Errorf("run golang build %s: %s", projectGo.Abs(), err)
		}
		projectGo.Manifest().Exec = build.Out
		if err = project.Dist(ctx, projectGo); err != nil {
			return
		}

		if slices.Contains(args, "pack") {
			if err = dist.Pack(projectGo.Dist_()); err != nil {
				return
			}
		}
	}
	return
}
