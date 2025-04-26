package rpc

import (
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"io"
)

func (r Rpc) Conn(target, query string) (c rpc.Conn, err error) {
	conn, err := r.Apphost.Query(target, query, nil)
	if err != nil {
		return
	}
	c = r.client(conn)
	return
}

func (r Rpc) client(conn io.ReadWriteCloser) *stream.Client {
	s := stream.Client{}
	s.Serializer = stream.NewSerializer(conn)
	return &s
}

func (r Rpc) serializer(conn io.ReadWriteCloser) *stream.Serializer {
	s := stream.NewSerializer(conn)
	s.Codec = r.codec()
	s.Logger(r.Log)
	return s
}
