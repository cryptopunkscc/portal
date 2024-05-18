package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"path"
)

func FindInPath[T target.Source](src string) (in <-chan T) {
	return target.Stream[T](Resolve, target.NewModule(src))
}

func FindInFS[T target.Source](src fs.FS) (in <-chan T) {
	return target.Stream[T](Resolve, target.NewModuleFS(src))
}

func Resolve(module target.Source) (result target.Source, err error) {
	if path.Base(module.Path()) == "node_modules" {
		return nil, fs.SkipDir
	}
	module = module.Lift()
	bundle, err := ResolveBundle(module)
	if err == nil {
		result = bundle
		return
	}
	if path.Ext(module.Path()) != "" && module.Path() != "." {
		err = nil
		return
	}
	nodeModule, err := ResolveNodeModule(module)
	if err == nil {
		if result, err = ResolvePortalNodeModule(nodeModule); err == nil {
			return
		}
		result = nodeModule
		err = nil
		return
	}
	result, err = ResolvePortalRawModule(module)
	if err == nil {
		err = fs.SkipDir
		return
	}
	err = nil
	return
}
