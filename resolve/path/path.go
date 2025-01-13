package path

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"io/fs"
)

func Resolver(source target.Source) target.Path {
	return func(port string) (path string, err error) {
		for _, t := range apps.Resolver[target.Bundle_]().List(source) {
			if t.Manifest().Match(port) {
				path = t.Abs()
				return
			}
		}
		err = fs.ErrNotExist
		return
	}
}
