package rpc

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/port"
	"io"
)

type Flow struct{ *Serializer }

func NewFlow(conn io.ReadWriteCloser) *Flow {
	s := Flow{&Serializer{}}
	s.setConn(conn)
	return &s
}

func QueryFlow(
	identity id.Identity,
	service string,
	path ...string,
) (s Conn, err error) {
	query, err := Apphost.Query(identity, port.Format(service, path...))
	if err != nil {
		return
	}
	return NewFlow(query), nil
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
