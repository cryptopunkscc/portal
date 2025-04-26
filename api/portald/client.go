package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"io"
)

func Client(rpc rpc.Rpc) Conn { return Conn{rpc.Request("portald")} }

type Conn struct{ rpc.Conn }

func (p Conn) Join()        { _ = rpc.Command(p, "join") }
func (p Conn) Ping() error  { return rpc.Command(p, "ping") }
func (p Conn) Close() error { return rpc.Command(p, "close") }
func (p Conn) Open(opt *OpenOpt, args ...string) error {
	var argv []any
	if opt != nil {
		argv = []any{opt}
	}
	for _, arg := range args {
		argv = append(argv, arg)
	}
	return rpc.Call(p.Copy(), "open", argv...)
}
func (p Conn) Connect(opt *OpenOpt, args ...string) (rwc io.ReadWriteCloser, err error) {
	var argv []any
	if opt != nil {
		argv = []any{opt}
	}
	for _, arg := range args {
		argv = append(argv, arg)
	}
	c := p.Copy()
	if err = rpc.Call(c, "connect", argv...); err != nil {
		return
	}
	rwc = c
	return
}
func (p Conn) Api() (cmd.Handlers, error) { return rpc.Query[cmd.Handlers](p, "api") }

type OpenOpt struct {
	Schema string `query:"s"`
	Order  []int  `query:"o"`
}

var portaldOpenOptKey = &OpenOpt{}

func (o *OpenOpt) Load(ctx context.Context) {
	if value, ok := ctx.Value(portaldOpenOptKey).(*OpenOpt); ok {
		*o = *value
	}
}

func (o *OpenOpt) Save(ctx *context.Context) {
	*ctx = context.WithValue(*ctx, portaldOpenOptKey, o)
}
