package exec

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

type source struct {
	executable target.Source
}

func (s *source) Executable() (t target.Source) {
	return s.executable
}