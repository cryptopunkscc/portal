package go_build

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/target"
	"os"
	"path/filepath"
	"runtime"
)

type runner struct{ platforms []string }

func Runner(platforms ...string) target.Run[target.ProjectGo] {
	return runner{platforms}.Run
}

func (g runner) Run(ctx context.Context, project target.ProjectGo) (err error) {
	log := plog.Get(ctx).Type(g).Set(&ctx)
	if err = deps.RequireBinary("go"); err != nil {
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

	if err = os.RemoveAll(filepath.Join(project.Abs(), "dist")); err != nil {
		log.W().Println(err)
	}
	for _, platform := range g.platforms {
		build, ok := project.Build()[platform]
		if !ok {
			build, ok = project.Build()["default"]
		}
		if ok {
			cmd = cmd.Parse(build.Cmd).AddEnv(build.Env...).AddEnv("GOOS=" + platform)
		}
		if err = cmd.Build().Run(); err != nil {
			return fmt.Errorf("run golang build %s: %s", project.Abs(), err)
		}
		project.Manifest().Exec = build.Exec
		if err = dist.Dist(ctx, project); err != nil {
			return
		}
	}
	return
}
