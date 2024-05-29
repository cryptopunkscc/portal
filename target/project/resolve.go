package project

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
	"github.com/cryptopunkscc/go-astral-js/target/npm"
	targetSource "github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
)

func FromPath(src string) (project target.Project, err error) {
	nodeModule, err := npm.ResolveNodeModule(targetSource.FromPath(src))
	if err != nil {
		return
	}
	return Resolve(nodeModule)
}

func Resolve(nodeModule target.NodeModule) (project target.Project, err error) {
	m := target.Manifest{}
	sub, err := fs.Sub(nodeModule.Files(), nodeModule.Path())
	if err != nil {
		return
	}
	if err = manifest.Load(&m, sub, target.PackageJsonFilename); err != nil {
		return
	}
	if err = manifest.Load(&m, sub, target.PortalJsonFilename); err != nil {
		return
	}
	project = &source{NodeModule: nodeModule, manifest: &m}
	switch {
	case project.Type().Is(target.TypeFrontend):
		project = &frontend{Project: project}
	case project.Type().Is(target.TypeBackend):
		project = &backend{Project: project}
	}
	return
}

func Dist[T target.Dist](project target.Project) (t T) {
	resolve := target.Any[T](target.Try(dist.Resolve))
	for _, t = range targetSource.List[T](resolve, project) {
		return
	}
	return
}
