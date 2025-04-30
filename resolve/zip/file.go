package zip

import (
	"github.com/cryptopunkscc/portal/api/target"
)

type File_ struct {
	target.Source
	file target.Source
}

var _ target.Bundle = File_{}

func (f File_) Package() target.Source { return f.file }
