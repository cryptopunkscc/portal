package apphost

import (
	"github.com/cryptopunkscc/portal/api/apphost"
)

func init() {
	apphost.Connect = Connect
	apphost.DefaultClient = Default
	apphost.DefaultCached = Cached(Default)
}
