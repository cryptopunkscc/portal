package target2

import (
	"io/fs"
	"reflect"
)

// Set all Source from a given dir.
func Set[T Base](resolve Resolve[T], from ...Source) (out []T) {
	out = List(resolve, from...)
	Reduce[T](&out)
	return
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
	_ = fs.WalkDir(from.Files(), from.Path(), func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipAll
		}
		m, err := from.Sub(src)
		if err != nil {
			return nil
		}
		s, err := resolve(m)
		if any(s) != nil && !reflect.ValueOf(s).IsNil() {
			out = append(out, s)
		}
		return err
	})
	return
}

// Reduce list removing duplicated elements by its Manifest.Package.
func Reduce[T Base](list *[]T) {
	var out []T
	l := *list
	m := make(map[string]T)
	for _, t := range l {
		if _, ok := m[t.Manifest().Package]; ok {
			continue
		}
		m[t.Manifest().Package] = t
		out = append(out, t)
	}
	*list = out
	return
}
