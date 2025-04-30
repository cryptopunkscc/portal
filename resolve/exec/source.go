package exec

import (
	"github.com/cryptopunkscc/portal/api/target"
)

type Source struct{ exec target.Source }

func (e Source) Executable() target.Source { return e.exec }
