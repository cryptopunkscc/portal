package build

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"runtime"
	"slices"
	"strings"
)

type runner struct {
	targets map[string][]string
}

func (r *runner) loadTargets(args ...string) (err error) {
	r.targets = make(map[string][]string)
	for _, arg := range args {
		a := strings.Split(arg, "-")
		os := a[0]
		r.targets[os] = append(r.targets[os], a[1:]...)
	}
	return
}

func (r *runner) Run(ctx context.Context, project target.Project_, args ...string) (err error) {
	if !project.Changed() && !slices.Contains(args, "clean") {
		return
	}

	if r.targets == nil {
		r.targets = map[string][]string{
			runtime.GOOS: {runtime.GOARCH},
		}
	}

	//for os, archs := range r.targets {
	//	for i, arch := range archs {
	//
	//	}
	//}
	//
	//getBuild(project)
	//
	//for s, build := range project.Build() {
	//	build
	//}
	//
	//cmd := exec.Cmd{}.Default()
	return
}

func getBuilds(project target.Project_) target.Builds {
	return nil
}

func getBuild(project target.Project_, args ...string) target.Build {
	builds := project.Build()
	var found []target.Build
	if build, ok := builds["default"]; ok {
		found = append(found, build)
	}
	if len(args) > 0 {
		if build, ok := builds[args[0]]; ok {
			found = append(found, build)
		}
	}
	if len(args) > 1 {
		if build, ok := builds[args[0]+"-"+args[1]]; ok {
			found = append(found, build)
		}
	}
	return target.MergeBuilds(found...)
}
