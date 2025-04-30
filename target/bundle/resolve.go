package bundle

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/zip"
)

var Resolve_ = Resolver(dist.Resolver(target.Resolve_))

func Resolver[T any](resolve target.Resolve[target.Dist[T]]) target.Resolve[target.AppBundle[T]] {
	return func(src target.Source) (bundle target.AppBundle[T], err error) {
		bundle, ok := src.(target.AppBundle[T])
		if ok {
			return
		}

		s := Source[T]{}
		if s.bundle, err = zip.Resolve(src); err != nil {
			return
		}
		if s.Dist, err = resolve(s.bundle); err != nil {
			return
		}
		bundle = s
		return
	}
}
