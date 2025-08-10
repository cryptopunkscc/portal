package project

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/dist"
)

var Resolve_ = Resolver(target.Resolve_)

func Resolver[T any](resolve target.Resolve[T]) target.Resolve[target.Project[T]] {
	return func(source target.Source) (project target.Project[T], err error) {
		defer plog.TraceErr(&err)
		if !source.IsDir() {
			err = target.ErrNotTarget
			return
		}

		if _, err = resolve(source); err != nil {
			return
		}
		s := &Source[T]{}
		if err = s.manifest.LoadFrom(source.FS()); err != nil {
			err = plog.Err(err, source.Abs())
			return
		}
		s.resolveDist = dist.Resolver(resolve)
		s.Source = source
		project = s
		return
	}
}
