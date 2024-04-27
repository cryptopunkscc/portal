package content

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/go-astral-js/pkg/android"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"io"
)

type Client struct {
	id.Identity
	rpc.Conn
}

func (c *Client) Connect() (err error) {
	conn, err := astral.Query(c.Identity, android.ContentPort)
	if err == nil {
		c.Conn = rpc.NewFlow(conn)
	}
	return
}

func (c *Client) Info(uri string) (files android.Info, err error) {
	if err = c.Connect(); err != nil {
		return
	}
	defer c.Close()
	return rpc.Query[android.Info](c.Conn, "info", uri)
}

func (c *Client) Reader(uri string, offset int64) (reader io.ReadCloser, err error) {
	if err = c.Connect(); err != nil {
		return
	}
	if err = rpc.Call(c.Conn, "reader", uri, offset); err != nil {
		return
	}
	reader = c.Conn
	return
}
