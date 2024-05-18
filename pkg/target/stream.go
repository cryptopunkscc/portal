package target

import (
	"io/fs"
	"reflect"
)

// Stream all portal targets in a given dir and stream through the returned channel.
// Possible types are: NodeModule, PortalNodeModule, PortalRawModule, Bundle,
func Stream[T Source](resolve Resolve, from Source) (in <-chan T) {
	out := make(chan T)
	in = out
	go func() {
		defer close(out)
		_ = fs.WalkDir(from.Files(), from.Path(), func(src string, d fs.DirEntry, err error) error {
			if err != nil {
				return fs.SkipAll
			}

			m := NewModuleFS(from.Files(), src, from.Abs())
			s, err := resolve(m)
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
