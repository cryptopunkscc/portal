package rpc

import (
	"github.com/cryptopunkscc/portal/pkg/util/rpc/cmd"
)

type Rpc interface {
	Router(handler cmd.Handler) Router
}
