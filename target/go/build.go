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
	"slices"
	"strings"
)

func BuildRunner(platforms ...string) *target.SourceRunner[target.ProjectGo] {
	return &target.SourceRunner[target.ProjectGo]{
		Resolve: ResolveProject,
		Runner:  BuildProject(platforms...),
	}
}

func BuildProject(platforms ...string) target.Run[target.ProjectGo] {
	p := make([][]string, len(platforms))
	for i, platform := range platforms {
		p[i] = strings.SplitN(platform, "/", 2)
	}
	return buildRunner{p}.Run
}

type buildRunner struct{ platforms [][]string }

func (g buildRunner) Run(ctx context.Context, projectGo target.ProjectGo, args ...string) (err error) {
	log := plog.Get(ctx).Type(g).Set(&ctx)
	if err = deps.RequireBinary("go"); err != nil {
		return
	}

	if !projectGo.Changed() && !slices.Contains(args, "clean") {
		return
	}

	log.Printf("go build %T %s %v", projectGo, projectGo.Abs(), g.platforms)

	if slices.Contains(args, "clean") {
		if err = os.RemoveAll(filepath.Join(projectGo.Abs(), "dist")); err != nil {
			log.W().Println(err)
		}
	}

	platforms := g.platforms
	if len(platforms) == 0 {
		platforms = [][]string{{}}
	}

	var cmd exec.Cmd
	for _, platform := range platforms {
		cmd, err = goBuild(projectGo, platform...)
		if err != nil {
			return
		}

		if err = cmd.Build().Run(); err != nil {
			return fmt.Errorf("run golang build %s: %s", projectGo.Abs(), err)
		}

		if err = project.Dist2(ctx, projectGo); err != nil {
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

func goBuild(projectGo target.ProjectGo, platform ...string) (cmd exec.Cmd, err error) {
	defer plog.TraceErr(&err)
	cmd = exec.Cmd{
		Cmd:  "go",
		Args: []string{"build", "-o", "dist/main"},
		Dir:  projectGo.Abs(),
	}.Default()

	b := projectGo.Build().Get(platform...)
	cmd, err = cmd.Parse(b.Cmd)
	if err != nil {
		return
	}

	if cmd, err = cmd.Parse(b.Cmd); err != nil {
		return
	}
	cmd = cmd.AddEnv(b.Env...).AddEnv("GOOS="+b.Target.OS, "GOARCH="+b.Target.Arch)
	if err = cmd.Build().Run(); err != nil {
		err = fmt.Errorf("run golang build %s: %s", projectGo.Abs(), err)
		return
	}
	return
}
