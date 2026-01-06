package apphost

import (
	"context"
	"io"

	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func (a *Adapter) PortaldLegacy() PortaldClientLegacy {
	return PortaldClientLegacy{a.Rpc().Request("portald")}
}

type PortaldClientLegacy struct{ rpc.Conn }

func (p PortaldClientLegacy) Join()        { _ = rpc.Command(p, "join") }
func (p PortaldClientLegacy) Ping() error  { return rpc.Command(p, "ping") }
func (p PortaldClientLegacy) Close() error { return rpc.Command(p, "close") }
func (p PortaldClientLegacy) Open(opt *OpenOptLegacy, args ...string) error {
	var argv []any
	if opt != nil {
		argv = []any{opt}
	}
	for _, arg := range args {
		argv = append(argv, arg)
	}
	return rpc.Call(p.Copy(), "open", argv...)
}
func (p PortaldClientLegacy) Connect(opt *OpenOptLegacy, args ...string) (rwc io.ReadWriteCloser, err error) {
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
func (p PortaldClientLegacy) Api() (cmd.Handlers, error) { return rpc.Query[cmd.Handlers](p, "api") }

type OpenOptLegacy struct {
	Schema string `query:"s"`
	Order  []int  `query:"o"`
}

var portaldOpenOptKeyLegacy = &OpenOptLegacy{}

func (o *OpenOptLegacy) Load(ctx context.Context) {
	if value, ok := ctx.Value(portaldOpenOptKeyLegacy).(*OpenOptLegacy); ok {
		*o = *value
	}
}

func (o *OpenOptLegacy) Save(ctx *context.Context) {
	*ctx = context.WithValue(*ctx, portaldOpenOptKeyLegacy, o)
}
