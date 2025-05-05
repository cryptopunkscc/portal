package exec

import (
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
)

type Source struct {
	executable target.Source
	target     manifest.Target
}

var _ target.Exec = Source{}

func (e Source) Executable() target.Source { return e.executable }
func (e Source) Target() manifest.Target   { return e.target }
