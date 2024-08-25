package rpc

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	"io"
)

type Request struct {
	*Serializer
	service string
	client  apphost.Client
}

func (r *Request) Client(client apphost.Client) *Request {
	r.client = client
	return r
}

func newRequest(identity id.Identity, service string) *Request {
	return &Request{
		Serializer: &Serializer{remoteID: identity},
		service:    service,
		client:     Apphost,
	}
}

func NewRequest(identity id.Identity, path ...string) *Request {
	return newRequest(identity, port.New(path...).String())
}

func (r *Request) Copy() Conn {
	rr := newRequest(r.remoteID, r.service)
	rr.client = r.client
	if r.logger != nil {
		rr.Logger(r.logger.Logger)
	}
	return rr
}

func (r *Request) Flush() {
	if r.WriteCloser != nil {
		_ = r.WriteCloser.Close()
	}
}

func (r *Request) Call(method string, value any) (err error) {
	// build base query
	query := ""
	switch {
	case r.service == "":
		query = method
	case method == "":
		query = r.service
	default:
		query = r.service + "." + method
	}

	// marshal args
	if value != nil {
		if query != "" {
			query += "?"
		}
		if r.marshal == nil {
			r.setupEncoding()
		}
		args, err := r.marshal(value)
		if err != nil {
			return plog.Err(err)
		}
		query += string(args)
	}

	// log query
	if r.logger != nil {
		r.logger.Println("~>", query)
	}

	// query stream
	var conn io.ReadWriteCloser
	if conn, err = r.client.Query(r.RemoteIdentity(), query); err != nil {
		return err
	}

	// setup
	r.setConn(conn)
	return
}
