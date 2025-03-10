package path

import (
	"github.com/cryptopunkscc/portal/api/target"
	"io/fs"
)

func Resolver[T target.Portal_](resolver target.Resolve[T], source target.Source) target.Path {
	return func(pkg string) (path string, err error) {
		for _, t := range resolver.List(source) {
			if t.Manifest().Match(pkg) {
				path = t.Abs()
				return
			}
		}
		err = fs.ErrNotExist
		return
	}
}
