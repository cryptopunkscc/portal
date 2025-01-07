package apphost

import (
	"encoding/json"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream/query"
	"io"
)

func RpcRequest(identity id.Identity, query string) rpc.Conn {
	return newRequest(Client, identity, query)
}

func (r RpcBase) Request(identity id.Identity, query string) rpc.Conn {
	return newRequest(r.client, identity, query)
}

func newRequest(client apphost.Client, identity id.Identity, q string) rpc.Conn {
	return &rpcRequest{
		Serializer: &stream.Serializer{
			MarshalArgs: query.Marshal,
			Marshal:     json.Marshal,
			Unmarshal:   json.Unmarshal,
		},
		query:    q,
		remoteID: identity,
		client:   client,
	}
}

type rpcRequest struct {
	*stream.Serializer
	query    string
	remoteID id.Identity
	logger   plog.Logger
	client   apphost.Client
}

func (r *rpcRequest) Logger(logger plog.Logger) {
	r.logger = logger
}

func (r *rpcRequest) Copy() rpc.Conn {
	rr := newRequest(r.client, r.remoteID, r.query)
	if r.logger != nil {
		rr.Logger(r.logger)
	}
	return rr
}

func (r *rpcRequest) Flush() {
	if r.Closer != nil {
		_ = r.Close()
	}
}

func (r *rpcRequest) Call(method string, value any) (err error) {
	// build base query
	q := ""
	switch {
	case r.query == "":
		q = method
	case method == "":
		q = r.query
	default:
		q = r.query + "." + method
	}

	// marshal args
	if value != nil {
		if q != "" {
			q += "?"
		}
		if r.MarshalArgs == nil {
			r.MarshalArgs = query.Marshal
		}
		args, err := r.MarshalArgs(value)
		if err != nil {
			return plog.Err(err)
		}
		q += string(args)
	}

	// log query
	if r.logger != nil {
		r.logger.Println("~>", q)
	}

	// query stream
	var conn io.ReadWriteCloser
	if conn, err = r.client.Query(r.remoteID, q); err != nil {
		return err
	}

	// setup
	if r.logger != nil {
		conn = rpc.NewConnLogger(conn, r.logger)
	}
	serializer := stream.NewSerializer(conn)
	serializer.MarshalArgs = r.MarshalArgs
	serializer.Marshal = r.Marshal
	serializer.Unmarshal = r.Unmarshal
	r.Serializer = serializer
	return
}
