package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"path"
	"reflect"
)

func FindInPath[T target.Source](src string) (in <-chan T) {
	return Find[T](NewModule(src))
}

func FindInFS[T target.Source](src fs.FS) (in <-chan T) {
	return Find[T](NewModuleFS(src, "."))
}

// Find all portal targets in a given dir and stream through the returned channel.
// Possible types are: NodeModule, PortalNodeModule, PortalRawModule, Bundle,
func Find[T target.Source](source target.Source) (in <-chan T) {
	out := make(chan T)
	in = out
	go func() {
		defer close(out)
		if source.Type().Is(target.Bundle) {
			var sources target.Source
			sources, _ = ResolveBundle(NewModule(source.Abs()))
			if sources != nil && !reflect.ValueOf(sources).IsNil() {
				switch t := sources.(type) {
				case T:
					out <- t
				}
			}
			return
		}

		_ = fs.WalkDir(source.Files(), source.Path(), func(src string, d fs.DirEntry, err error) error {
			if err != nil {
				return fs.SkipAll
			}
			s, err := resolveSources(source, src, d)
			if s != nil && !reflect.ValueOf(s).IsNil() {
				switch t := s.(type) {
				case T:
					out <- t
				}
			}
			return err
		})
	}()
	return
}

func resolveSources(root target.Source, src string, d fs.DirEntry) (result target.Source, err error) {
	if d.Name() == "node_modules" {
		return nil, fs.SkipDir
	}
	module := NewModuleFS(root.Files(), src)
	module.abs = path.Join(root.Abs(), src)
	bundle, err := ResolveBundle(module)
	if err == nil {
		result = bundle
		return
	}
	if d.Type().IsRegular() {
		err = nil
		return
	}
	module = module.Lift()
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
