package rpc

import "github.com/cryptopunkscc/portal/pkg/rpc/cmd"

type Rpc interface {
	Format(name string) Rpc
	Conn(target, query string) (Conn, error)
	Request(target string, query ...string) Conn
	Router(handler cmd.Handler) Router
}
