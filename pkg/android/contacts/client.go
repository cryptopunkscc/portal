package contacts

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/runtime/rpc2"
)

type Client struct {
	rpc.Conn
}

func (c Client) Connect(identity *astral.Identity, port string) (client Client, err error) {
	client.Conn, err = rpc.QueryFlow(identity, port)
	return
}

func (c Client) Contacts() (<-chan []Contact, error) {
	return rpc.Subscribe[[]Contact](c, "contacts")
}
