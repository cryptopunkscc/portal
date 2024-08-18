package main

import "github.com/cryptopunkscc/portal/runtime/rpc"

type portalApi struct{ rpc.Conn }

func (p portalApi) Await()                { _ = rpc.Command(p, "") }
func (p portalApi) Ping() error           { return rpc.Command(p, "ping") }
func (p portalApi) Open(src string) error { return rpc.Command(p, "open", src) }
func (p portalApi) Close() error          { return rpc.Command(p, "close") }
