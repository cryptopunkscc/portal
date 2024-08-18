package client

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

func Portal(pkg string) portal.Client { return PortalClient{rpc.NewRequest(id.Anyone, pkg)} }

type PortalClient struct{ rpc.Conn }

func (p PortalClient) Await()                { _ = rpc.Command(p, "") }
func (p PortalClient) Ping() error           { return rpc.Command(p, "ping") }
func (p PortalClient) Open(src string) error { return rpc.Command(p, "open", src) }
func (p PortalClient) Close() error          { return rpc.Command(p, "close") }
