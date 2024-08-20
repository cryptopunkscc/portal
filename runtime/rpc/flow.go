package rpc

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/port"
	"io"
)

type Flow struct {
	*Serializer
	client apphost.Client
}

func QueryFlow(identity id.Identity, service string, path ...string) (Conn, error) {
	return NewFlow().Query(identity, service, path...)
}

func NewFlow() *Flow { return &Flow{&Serializer{}, Apphost} }

func (conn *Flow) Client(client apphost.Client) *Flow {
	conn.client = client
	return conn
}

func (conn *Flow) Conn(client io.ReadWriteCloser) *Flow {
	conn.setConn(client)
	return conn
}

func (conn *Flow) Query(
	identity id.Identity,
	service string,
	path ...string,
) (Conn, error) {
	query, err := conn.client.Query(identity, port.Format(service, path...))
	if err != nil {
		return nil, err
	}
	return conn.Conn(query), nil
}

func (conn *Flow) Call(method string, value any) (err error) {
	query := []byte(method)
	if value != nil {
		var bytes []byte
		if bytes, err = conn.marshal(value); err != nil {
			return
		}
		query = append(query, bytes...)
	}
	query = append(query, []byte("\n")...)
	writer := conn.WriteCloser
	if conn.logger != nil {
		writer = conn.logger
	}
	_, err = writer.Write(query)
	return
}

func (conn *Flow) Copy() Conn {
	return conn
}

func (conn *Flow) Flush() {
	// no-op
}
