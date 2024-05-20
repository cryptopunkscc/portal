package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/manifest"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/npm"
	"io/fs"
)

func New(src string) (module target.Project, err error) {
	nodeModule, err := npm.ResolveNodeModule(target.NewModule(src))
	if err != nil {
		return
	}
	return Resolve(nodeModule)
}

func Resolve(t target.NodeModule) (b target.Project, err error) {
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
	b = &source{NodeModule: t, manifest: &m}
	switch {
	case b.Type().Is(target.TypeFrontend):
		b = &frontend{Project: b}
	case b.Type().Is(target.TypeBackend):
		b = &backend{Project: b}
	}
	return
}
