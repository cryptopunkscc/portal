package rpc

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Rpc struct {
	Apphost apphost.Client
	Log     plog.Logger
}
