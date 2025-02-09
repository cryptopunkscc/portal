package apphost

import (
	"encoding/json"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream/query"
	"io"
	"log"
)

func (r RpcBase) Request(target string, query ...string) rpc.Conn {
	return newRequest(r.client, target, query)
}

func newRequest(client apphost.Client, target string, q []string) *rpcRequest {
	return &rpcRequest{
		Serializer: &stream.Serializer{
			MarshalArgs: query.Marshal,
			Marshal:     json.Marshal,
			Unmarshal:   json.Unmarshal,
		},
		client: client,
		target: target,
		query:  q,
	}
}

type rpcRequest struct {
	*stream.Serializer
	logger   plog.Logger
	client   apphost.Client
	target   string
	targetID *astral.Identity
	query    []string
}

func (r *rpcRequest) Logger(logger plog.Logger) {
	r.logger = logger
}

func (r *rpcRequest) Copy() rpc.Conn {
	rr := newRequest(r.client, r.target, r.query)
	rr.targetID = r.targetID
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
	if r.client == nil {
		r.client = apphost.DefaultClient
	}
	// build base query
	p := port.New(r.query...)
	if method != "" {
		p = p.Add(method)
	}
	q := p.String()

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

	if r.targetID == nil {
		r.targetID, err = r.client.Resolve(r.target)
		if err != nil {
			log.Println("failed to resolve target", r.target, err)
			return plog.Err(err)
		}
	}

	// query stream
	var conn io.ReadWriteCloser
	if conn, err = r.client.Query(r.targetID.String(), q, nil); err != nil {
		return err
	}

	// setup
	serializer := stream.NewSerializer(conn)
	serializer.MarshalArgs = r.MarshalArgs
	serializer.Marshal = r.Marshal
	serializer.Unmarshal = r.Unmarshal
	if r.logger != nil {
		serializer.Logger(r.logger)
	}
	r.Serializer = serializer
	return
}
