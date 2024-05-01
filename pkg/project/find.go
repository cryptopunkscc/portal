package project

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"path"
)

func RawTargets(files fs.FS) <-chan PortalRawModule {
	return Find[PortalRawModule](files, ".")
}

func DevTargets(files fs.FS) <-chan PortalNodeModule {
	return Find[PortalNodeModule](files, ".")
}

func Bundles(files fs.FS, dir string) <-chan Bundle {
	return Find[Bundle](files, dir)
}

func Apps(files fs.FS) <-chan target.App {
	return Find[target.App](files, ".")
}

func Find[T target.Source](files fs.FS, dir string) (in <-chan T) {
	out := make(chan T)
	in = out
	var t T
	go func() {
		defer close(out)
		_ = fs.WalkDir(files, dir, func(src string, d fs.DirEntry, err error) error {
			if err != nil {
				return fs.SkipAll
			}
			sources, err := resolveSources(files, src, d, (any)(t) == nil)
			if sources != nil {
				switch t := sources.(type) {
				case T:
					out <- t
				}
			}
			return err
		})
	}()
	return
}

func resolveSources(files fs.FS, src string, d fs.DirEntry, pointer bool) (result target.Source, err error) {
	if d.Name() == "node_modules" {
		return nil, fs.SkipDir
	}
	if path.Ext(src) == ".portal" && d.Type().IsRegular() {
		var sub fs.FS
		if sub, err = fs.Sub(files, path.Dir(src)); err != nil {
			err = fmt.Errorf("fs.Sub: %v", err)
			return
		}
		var bundle *Bundle
		if bundle, err = newModuleFS(src, sub).Bundle(); err != nil {
			err = fmt.Errorf("newModuleFS: %v", err)
			return
		}
		err = nil
		if pointer {
			result = bundle
			return
		}
		result = *bundle
		return
	}

	sub, err := fs.Sub(files, src)
	if err != nil {
		return
	}
	module := newModuleFS(src, sub)
	nodeModule, err := module.NodeModule()
	if err == nil {
		var portalModule *PortalNodeModule
		portalModule, err = nodeModule.PortalNodeModule()
		if err == nil {
			if pointer {
				result = portalModule
				return
			}
			result = *portalModule
			return
		}
		err = nil
		if pointer {
			result = nodeModule
			return
		}
		result = *nodeModule
		return
	}

	rawModule, err := module.PortalRawModule()
	if err == nil {
		err = fs.SkipDir
		if pointer {
			result = rawModule
			return
		}
		result = *rawModule
		return
	}
	err = nil
	return
}
