package dist

import (
	"github.com/cryptopunkscc/portal/api/target"
)

type Source[T any] struct {
	target.Dist_
	target T
}

var _ target.Dist[any] = &Source[any]{}

func (s Source[T]) Runtime() T { return s.target }
