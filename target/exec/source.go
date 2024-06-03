package exec

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

type source struct {
	executable target.Source
}

func (s *source) Exec() (t target.Source) {
	return s.executable
}
