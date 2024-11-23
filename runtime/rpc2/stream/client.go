package stream

import rpc "github.com/cryptopunkscc/portal/runtime/rpc2"

type Client struct{ *Serializer }

func (conn *Client) Call(method string, value any) (err error) {
	query := []byte(method)
	if value != nil {
		var bytes []byte
		if bytes, err = conn.Marshal(value); err != nil {
			return
		}
		query = append(query, bytes...)
	}
	query = append(query, []byte("\n")...)
	_, err = conn.Write(query)
	return
}

func (conn *Client) Copy() rpc.Conn {
	return conn
}

func (conn *Client) Flush() {
	// no-op
}
