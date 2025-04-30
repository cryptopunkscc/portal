package project

import (
	"github.com/cryptopunkscc/portal/api/target"
)

var Resolve_ = Resolver(target.Resolve_)

func Resolver[T any](resolve target.Resolve[T]) target.Resolve[target.Project[T]] {
	return func(source target.Source) (project target.Project[T], err error) {
		return New(source, resolve)
	}
}
