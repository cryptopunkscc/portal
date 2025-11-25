package stream

import (
	"bytes"

	"github.com/cryptopunkscc/portal/pkg/rpc"
)

type Client struct{ *Serializer }

func (conn *Client) Call(method string, value any) (err error) {
	query := bytes.NewBufferString(method)
	if value != nil {
		var args []byte
		if args, err = conn.MarshalArgs(value); err != nil {
			return
		}
		query.WriteByte('?')
		query.Write(args)
	}
	query.WriteByte('\n')
	_, err = query.WriteTo(conn)
	return
}

func (conn *Client) Copy() rpc.Conn {
	return conn
}

func (conn *Client) Flush() {
	// no-op
}
