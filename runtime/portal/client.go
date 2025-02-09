package portal

import (
	apphostApi "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"io"
)

func NewClient(apphostClient apphostApi.Client) portal.Client {
	return ClientRpc{Conn: apphost.Rpc(apphostClient).Request("portal", "portal")}
}

type ClientRpc struct{ rpc.Conn }

func (p ClientRpc) Join() { _ = rpc.Command(p, "") }

func (p ClientRpc) Ping() error  { return rpc.Command(p, "ping") }
func (p ClientRpc) Close() error { return rpc.Command(p, "close") }
func (p ClientRpc) Open(opt *portal.OpenOpt, args ...string) error {
	var argv []any
	if opt != nil && opt.Schema != "" {
		argv = []any{*opt}
	}
	for _, arg := range args {
		argv = append(argv, arg)
	}
	return rpc.Call(p.Copy(), "open", argv...)
}
func (p ClientRpc) Connect(opt *portal.OpenOpt, args ...string) (rwc io.ReadWriteCloser, err error) {
	var argv []any
	if opt != nil && opt.Schema != "" {
		argv = []any{opt}
	}
	for _, arg := range args {
		argv = append(argv, arg)
	}
	conn := p.Copy()
	if err = rpc.Call(conn, "connect", argv...); err != nil {
		return
	}
	rwc = conn
	return
}
