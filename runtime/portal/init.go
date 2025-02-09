package portal

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
)

func init() {
	portal.DefaultClient = NewClient(apphost.DefaultClient)
}
