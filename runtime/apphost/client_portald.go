package apphost

import (
	"github.com/cryptopunkscc/portal/runtime/rpc"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
	"io"
)

func (a *Adapter) Portald() Portald { return Portald{Default.Rpc().Request("portal")} }

type Portald struct{ rpc.Conn }

func (p Portald) Join()        { _ = rpc.Command(p, "join") }
func (p Portald) Ping() error  { return rpc.Command(p, "ping") }
func (p Portald) Close() error { return rpc.Command(p, "close") }
func (p Portald) Open(opt *PortaldOpenOpt, args ...string) error {
	var argv []any
	if opt != nil {
		argv = []any{opt}
	}
	for _, arg := range args {
		argv = append(argv, arg)
	}
	return rpc.Call(p.Copy(), "open", argv...)
}
func (p Portald) Connect(opt *PortaldOpenOpt, args ...string) (rwc io.ReadWriteCloser, err error) {
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
func (p Portald) Api() (cmd.Handlers, error) { return rpc.Query[cmd.Handlers](p, "api") }

type PortaldOpenOpt struct {
	Schema string `query:"s"`
	Order  []int  `query:"o"`
}
