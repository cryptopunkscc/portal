package rpc

import "github.com/cryptopunkscc/portal/pkg/rpc/cmd"

type Rpc interface {
	Router(handler cmd.Handler) Router
}
