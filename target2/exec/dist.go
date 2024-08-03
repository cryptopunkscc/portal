package exec

import (
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/dist"
)

var ResolveDist = dist.Resolver[Exec](ResolveExec)
