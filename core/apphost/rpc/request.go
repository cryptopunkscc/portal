package rpc

import (
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
)

func (r *Rpc) Request(target string, query ...string) rpc.Conn {
	rr := &rpcRequest{
		Rpc:    *r,
		target: target,
		query:  query,
	}
	rr.Serializer.Codec = r.codec()
	return rr
}

type rpcRequest struct {
	Rpc
	stream.Serializer
	target   string
	targetID *astral.Identity
	query    []string
}

func (r *rpcRequest) Copy() rpc.Conn {
	rr := *r
	rr.Serializer = stream.Serializer{}
	rr.Serializer.Codec = r.Serializer.Codec
	return &rr
}

func (r *rpcRequest) Logger(logger plog.Logger) {
	r.Log = logger
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
		var args []byte
		if args, err = r.Serializer.MarshalArgs(value); err != nil {
			return
		}
		q += string(args)
	}

	// log query
	if r.Log != nil {
		r.Log.Println("~>", q)
	}

	if r.targetID == nil {
		if r.targetID, err = r.Apphost.Resolve(r.target); err != nil {
			return
		}
	}

	// query stream
	var conn io.ReadWriteCloser
	if conn, err = r.Apphost.Query(r.targetID.String(), q, nil); err != nil {
		return err
	}

	// setup
	s := stream.NewSerializer(conn)
	s.Codec = r.Serializer.Codec
	s.Logger(r.Log)
	r.Serializer = *s
	return
}
