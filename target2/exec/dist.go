package exec

import (
	. "github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/dist"
)

var ResolveDist = dist.Resolver[Exec](ResolveExec)
