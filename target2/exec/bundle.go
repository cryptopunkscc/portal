package exec

import (
	. "github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/bundle"
)

var ResolveBundle = bundle.Resolver[Exec](ResolveDist)
