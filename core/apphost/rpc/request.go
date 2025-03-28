package rpc

import (
	"encoding/json"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream/query"
	"io"
)

func (r Rpc) Request(target string, query ...string) rpc.Conn {
	return &rpcRequest{
		Serializer: newSerializer(),
		logger:     r.Log,
		client:     r.Apphost,
		target:     target,
		query:      query,
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

func (r *rpcRequest) Copy() rpc.Conn {
	return &rpcRequest{
		Serializer: newSerializer(),
		client:     r.client,
		target:     r.target,
		query:      r.query,
		targetID:   r.targetID,
		logger:     r.logger,
	}
}

func newSerializer() *stream.Serializer {
	return &stream.Serializer{
		MarshalArgs: query.Marshal,
		Marshal:     json.Marshal,
		Unmarshal:   json.Unmarshal,
		Ending:      []byte("\n"),
	}
}

func (r *rpcRequest) Logger(logger plog.Logger) {
	r.logger = logger
}

func (r *rpcRequest) Flush() {
	if r.Closer != nil {
		_ = r.Close()
	}
}

func (r *rpcRequest) Call(method string, value any) (err error) {
	defer plog.TraceErr(&err)
	// build base query
	p := apphost.NewPort(r.query...)
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
		var args []byte
		if args, err = r.MarshalArgs(value); err != nil {
			return
		}
		q += string(args)
	}

	// log query
	if r.logger != nil {
		r.logger.Println("~>", q)
	}

	if r.targetID == nil {
		if r.targetID, err = r.client.Resolve(r.target); err != nil {
			return
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
