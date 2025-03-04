package portald

import (
	"github.com/cryptopunkscc/portal/runtime/rpc"
	"github.com/cryptopunkscc/portal/runtime/rpc/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
	"io"
)

func NewClient() Client { return Client{apphost.Request("portal")} }

type Client struct{ rpc.Conn }

func (p Client) Join()        { _ = rpc.Command(p, "join") }
func (p Client) Ping() error  { return rpc.Command(p, "ping") }
func (p Client) Close() error { return rpc.Command(p, "close") }
func (p Client) Open(opt *OpenOpt, args ...string) error {
	var argv []any
	if opt != nil && opt.Schema != "" {
		argv = []any{*opt}
	}
	for _, arg := range args {
		argv = append(argv, arg)
	}
	return rpc.Call(p.Copy(), "open", argv...)
}
func (p Client) Connect(opt *OpenOpt, args ...string) (rwc io.ReadWriteCloser, err error) {
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
func (p Client) Api() (cmd.Handlers, error) { return rpc.Query[cmd.Handlers](p, "api") }

type OpenOpt struct {
	Schema string `query:"s"`
	Order  []int  `query:"o"`
}
