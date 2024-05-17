package source

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"io/fs"
	"reflect"
)

// Stream all portal targets in a given dir and stream through the returned channel.
// Possible types are: NodeModule, PortalNodeModule, PortalRawModule, Bundle,
func Stream[T target.Source](resolve target.Resolve[T], from target.Source) (in <-chan T) {
	out := make(chan T)
	in = out
	go func() {
		defer close(out)
		_ = fs.WalkDir(from.Files(), from.Path(), func(src string, d fs.DirEntry, err error) error {
			if err != nil {
				return fs.SkipAll
			}

			m := Resolve(from.Files(), src, from.Abs())
			s, err := resolve(m)
			if any(s) != nil && !reflect.ValueOf(s).IsNil() {
				out <- s
			}
			return err
		})
	}()
	return
}
