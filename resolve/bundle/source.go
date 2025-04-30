package bundle

import (
	"github.com/cryptopunkscc/portal/api/target"
)

type Source[T any] struct {
	target.Dist[T]
	bundle target.Bundle
}

func (s Source[T]) Package() target.Source { return s.bundle.Package() }
