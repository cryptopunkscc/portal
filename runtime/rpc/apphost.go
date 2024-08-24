package rpc

import (
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/portal/runtime/apphost"
)

var Apphost = apphost.Adapter(astral.Client)
