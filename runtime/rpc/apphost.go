package rpc

import (
	api "github.com/cryptopunkscc/portal/api/apphost"
	runtime "github.com/cryptopunkscc/portal/runtime/apphost"
)

var Apphost api.Client = runtime.Default()
