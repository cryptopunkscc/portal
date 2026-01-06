package apphost

import (
	"context"
	"io"

	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func (a *Adapter) Portald() PortaldClient {
	return PortaldClient{a}
}

type PortaldClient struct{ *Adapter }

func (c PortaldClient) Join() {
	_ = Call(nil, *c.Client, "portald.join", nil)
}

func (c PortaldClient) Ping() error {
	return Call(nil, *c.Client, "portald.join", nil)
}

func (c PortaldClient) Close() error {
	return Call(nil, *c.Client, "portald.close", nil)
}

func (c PortaldClient) Api() (rwc cmd.Handlers, err error) {
	return Receive[cmd.Handlers](nil, *c.Client, "portald.api", nil)
}

func (c PortaldClient) Connect(opt *OpenOpt) (rwc io.ReadWriteCloser, err error) {
	return c.Client.Query(nil, "portald.connect", opt)
}

func (c PortaldClient) Open(opt *OpenOpt) (err error) {
	return Call(nil, *c.Client, "portald.connect", opt)
}

type OpenOpt struct {
	Schema string `query:"s"`
	Order  string `query:"o"`
	Args   string `query:"a"`
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
