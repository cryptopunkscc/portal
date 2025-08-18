package golang

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/project"
	"os"
	"path/filepath"
	"runtime"
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
	if len(p) == 0 {
		p = [][]string{{runtime.GOOS, runtime.GOARCH}}
	}
	return buildRunner{p}.Run
}

type buildRunner struct{ platforms [][]string }

func (g buildRunner) Run(ctx context.Context, projectGo target.ProjectGo, args ...string) (err error) {
	log := plog.Get(ctx).Type(g).Set(&ctx)
	if err = deps.RequireBinary("go"); err != nil {
		return
	}

	clean := target.Op(&args, "clean")
	pack := target.Op(&args, "pack")

	log.Printf("go build %T %s %v", projectGo, projectGo.Abs(), g.platforms)

	for _, platform := range g.getPlatforms(&args) {
		b := projectGo.Build().Get(platform...)
		p := project.DistPath(b.Target)

		if clean || Changed(projectGo, platform...) {

			distPath := append([]string{projectGo.Abs(), "dist"}, p...)
			if err = os.RemoveAll(filepath.Join(distPath...)); err != nil {
				log.W().Println(err)
			}

			if err = goBuild(b, projectGo.Abs()); err != nil {
				return fmt.Errorf("run golang build %s: %s", projectGo.Abs(), err)
			}

			print(fmt.Sprintln("target: ", b.Target, projectGo.Abs()))
			if err = project.Dist(ctx, projectGo, b.Target); err != nil {
				return
			}
		}

		if pack {
			d := projectGo.Dist_(p...)
			abs := projectGo.Abs()
			if len(args) > 0 {
				abs = args[0]
			}
			if err = dist.Pack(d, abs); err != nil {
				return
			}
		}
	}
	return
}

func (g buildRunner) getPlatforms(args *[]string) (platforms [][]string) {
	goos := target.OpVal(args, "goos=")
	goarch := target.OpVal(args, "goarch=")
	if goos != "" {
		if goarch == "" {
			goarch = runtime.GOARCH
		}
		return [][]string{{goos, goarch}}
	}

	return g.platforms
}

func goBuild(build manifest.Build, abs string) (err error) {
	defer plog.TraceErr(&err)
	t := build.Target
	o := filepath.Join("dist", t.OS, t.Arch, "main")
	cmd := exec.Cmd{
		Cmd:  "go",
		Args: []string{"build", "-o", o},
		Dir:  abs,
	}.Default()

	build.Cmd = strings.ReplaceAll(build.Cmd, "$OUT", o)
	if cmd, err = cmd.Parse(build.Cmd); err != nil {
		return
	}
	cmd = cmd.AddEnv(build.Env...).AddEnv("GOOS="+t.OS, "GOARCH="+t.Arch)
	if err = cmd.Build().Run(); err != nil {
		err = fmt.Errorf("run golang build %s: %s", abs, err)
		return
	}
	return
}
