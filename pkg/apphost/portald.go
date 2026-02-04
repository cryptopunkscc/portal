package apphost

import (
	"context"
	"io"

	"github.com/cryptopunkscc/portal/pkg/util/rpc/cmd"
)

func (a *Adapter) Portald() PortaldClient {
	return PortaldClient{a}
}

type PortaldClient struct{ *Adapter }

func (c PortaldClient) Join() {
	//_ = Call(nil, *c.Client, "portald.join", nil)
}

func (c PortaldClient) Ping() error {
	//return Call(nil, *c.Client, "portald.join", nil)
	panic("not implemented")
}

func (c PortaldClient) Close() error {
	//return Call(nil, *c.Client, "portald.close", nil)
	panic("not implemented")
}

func (c PortaldClient) Api() (rwc cmd.Handlers, err error) {
	//return Receive[cmd.Handlers](nil, *c.Client, "portald.api", nil)
	panic("not implemented")

}

func (c PortaldClient) Connect(opt OpenOpt) (rwc io.ReadWriteCloser, err error) {
	//return c.Client.Query(nil, "portald.connect", opt)
	panic("not implemented")
}

func (c PortaldClient) Open(opt OpenOpt) (err error) {
	//return Call(nil, *c.Client, "portald.connect", opt)
	panic("not implemented")
}

type OpenOpt struct {
	App  string `query:"app"`
	Args string `query:"args"`
}

var portaldOpenOptKey = &OpenOpt{}

func (o *OpenOpt) Load(ctx context.Context) {
	//if value, ok := ctx.Value(portaldOpenOptKey).(*OpenOpt); ok {
	//	*o = *value
	//}
}

func (o *OpenOpt) Save(ctx *context.Context) {
	//*ctx = context.WithValue(*ctx, portaldOpenOptKey, o)
}
