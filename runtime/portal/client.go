package portal

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
)

func Client(pkg string) portal.Client { return ClientRpc{apphost.RpcRequest(id.Anyone, pkg)} }

type ClientRpc struct{ rpc.Conn }

func (p ClientRpc) Join()                 { _ = rpc.Command(p, "") }
func (p ClientRpc) Ping() error           { return rpc.Command(p, "ping") }
func (p ClientRpc) Open(src string) error { return rpc.Command(p, "open", src) }
func (p ClientRpc) Close() error          { return rpc.Command(p, "close") }
