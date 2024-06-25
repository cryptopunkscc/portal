package source

import (
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
	"reflect"
)

// List all target.Source from a given dir.
func List[T target.Source](resolve target.Resolve[T], from target.Source) (out []T) {
	out = []T{}
	//log.Println("starting stream", from.Abs(), from.Files(), reflect.TypeOf(from.Files()))
	_ = fs.WalkDir(from.Files(), from.Path(), func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipAll
		}
		m := FromFS(from.Files(), src, from.Abs())
		s, err := resolve(m)
		if any(s) != nil && !reflect.ValueOf(s).IsNil() {
			out = append(out, s)
		}
		return err
	})
	return
}
