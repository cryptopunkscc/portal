package project

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
	"github.com/cryptopunkscc/go-astral-js/target/npm"
	targetSource "github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
)

func FromPath(src string) (project target.ProjectNpm, err error) {
	nodeModule, err := npm.ResolveNodeModule(targetSource.FromPath(src))
	if err != nil {
		return
	}
	return ResolveNpm(nodeModule)
}

func ResolveNpm(t target.NodeModule) (project target.ProjectNpm, err error) {
	m := target.Manifest{}
	sub, err := fs.Sub(t.Files(), t.Path())
	if err != nil {
		return
	}
	if err = manifest.Load(&m, sub, target.PackageJsonFilename); err != nil {
		return
	}
	if err = manifest.Load(&m, sub, target.PortalJsonFilename); err != nil {
		return
	}
	src := portal{manifest: &m, Source: t}
	project = &nodeModule{NodeModule: t, portal: src}
	switch {
	case project.Type().Is(target.TypeFrontend):
		project = &html{ProjectNpm: project}
	case project.Type().Is(target.TypeBackend):
		project = &js{ProjectNpm: project}
	}
	return
}

func ResolveGo(t target.Source) (project target.ProjectGo, err error) {
	sub, err := fs.Sub(t.Files(), t.Path())
	if err != nil {
		return
	}
	m, err := manifest.Read(sub)
	if err != nil {
		return
	}
	mainGo, err := sub.Open("main.go")
	if err != nil {
		return
	}
	_ = mainGo.Close()
	src := portal{manifest: &m, Source: t.Lift()}
	project = &golang{portal: src}
	return
}

func Dist[T target.Dist](project target.Project) (t T) {
	resolve := target.Any[T](target.Try(dist.Resolve))
	for _, t = range targetSource.List[T](resolve, project) {
		return
	}
	return
}
