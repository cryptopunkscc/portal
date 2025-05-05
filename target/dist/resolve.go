package dist

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

var Resolve_ target.Resolve[target.Dist_] = resolve_

func resolve_(source target.Source) (portal target.Dist_, err error) {
	defer plog.TraceErr(&err)

	portal, ok := source.(target.Dist_)
	if ok {
		return
	}

	s := &Source_{}
	s.Source = source
	if err = s.manifest.LoadFrom(source.FS()); err != nil {
		return
	}
	portal = s
	return
}

func Resolver[T any](resolveT target.Resolve[T]) target.Resolve[target.Dist[T]] {
	return func(source target.Source) (portal target.Dist[T], err error) {
		defer plog.TraceErr(&err)

		portal, ok := source.(target.Dist[T])
		if ok {
			return
		}

		p, err := resolve_(source)
		if err != nil {
			return
		}

		t, err := resolveT(p)
		if err != nil {
			return
		}

		portal = &Source[T]{
			Dist_:  p,
			target: t,
		}
		return
	}
}
