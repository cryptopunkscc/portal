package apphost

import (
	"github.com/cryptopunkscc/astrald/lib/astral"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/runtime/apphost"
)

var Client api.Client = apphost.Cached(apphost.Adapter(astral.Client))
