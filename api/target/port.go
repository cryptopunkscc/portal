package target

import (
	"github.com/cryptopunkscc/portal/api/apphost"
)

var PortPortal = apphost.NewPort("portal")
var PortOpen = PortPortal.Add("open")
var PortMsg = PortPortal.Add("broadcast")
