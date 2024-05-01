package project

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"log"
	"path"
	"reflect"
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
	go func() {
		defer close(out)
		_ = fs.WalkDir(files, dir, func(src string, d fs.DirEntry, err error) error {
			if err != nil {
				return fs.SkipAll
			}
			sources, err := resolveSources(files, src, d)
			if sources != nil {
				t, ok := sources.(T)
				if ok {
					log.Println("load: ", reflect.TypeOf(t), src)
					out <- t
				}
			}
			return err
		})
	}()
	return
}

func resolveSources(files fs.FS, src string, d fs.DirEntry) (result any, err error) {
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
		bundle, err = newModuleFS(src, sub).Bundle()
		if err != nil {
			err = fmt.Errorf("newModuleFS: %v", err)
			return
		}
		result = *bundle
		err = nil
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
			result = *portalModule
			err = fs.SkipDir
			return
		}
		result = *nodeModule
		err = nil
		return
	}

	rawModule, err := module.PortalRawModule()
	if err == nil {
		result = *rawModule
		err = fs.SkipDir
		return
	}
	err = nil
	return
}
