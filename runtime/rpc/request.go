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

func (r *Request) Client(client apphost.Client) {
	r.client = client
}

func newRequest(identity id.Identity, service string) Conn {
	return &Request{
		Serializer: &Serializer{remoteID: identity},
		service:    service,
		client:     Apphost,
	}
}

func NewRequest(identity id.Identity, service string, path ...string) Conn {
	return newRequest(identity, port.Format(service, path...))
}

func (r *Request) Copy() Conn {
	rr := newRequest(r.remoteID, r.service)
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
