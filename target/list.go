package target

import (
	"errors"
	"io/fs"
	"log"
	"reflect"
	"slices"
)

// Set all Source from a given dir.
func Set[T Portal_](resolve Resolve[T], from ...Source) (out Portals[T]) {
	return Portals[T](List(resolve, from...)).Reduced()
}

// List all Source from a given dir.
func List[T any](resolve Resolve[T], from ...Source) (out []T) {
	for _, src := range from {
		o := resolve.list(src)
		out = append(out, o...)
	}
	return
}

// list all Source from a given dir.
func (resolve Resolve[T]) list(from Source) (out []T) {
	if !from.IsDir() {
		if t, err := resolve(from); err == nil {
			return append(out, t)
		}
	}
	_ = fs.WalkDir(from.Files(), from.Path(), func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println("Resolve list", err)
			return err
		}
		m, err := from.Sub(src)
		if err != nil {
			log.Println("Resolve sub", err)
			return nil
		}
		s, err := resolve(m)
		if errors.Is(err, fs.SkipDir) || errors.Is(err, fs.SkipAll) {
			return err
		}
		if any(s) != nil && !reflect.ValueOf(s).IsNil() {
			out = append(out, s)
		}
		return nil
	})
	return
}

func (p Portals[T]) SortBy(priority Priority) {
	slices.SortFunc(p, func(a, b T) int {
		return priority.Get(a) - priority.Get(b)
	})
}

func (p Portals[T]) Reduced() (reduced Portals[T]) {
	mem := make(map[string]T)
	for _, t := range p {
		if _, ok := mem[t.Manifest().Package]; ok {
			continue
		}
		mem[t.Manifest().Package] = t
		reduced = append(reduced, t)
	}
	return
}
