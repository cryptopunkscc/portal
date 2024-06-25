package exec

import (
	"github.com/cryptopunkscc/portal/target"
)

type source struct {
	executable target.Source
}

func (s *source) Executable() (t target.Source) {
	return s.executable
}
