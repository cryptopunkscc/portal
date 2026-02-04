package rpc

import (
	"io"

	"github.com/cryptopunkscc/portal/pkg/util/rpc/stream"
)

func (r *Rpc) client(conn io.ReadWriteCloser) *stream.Client {
	s := stream.Client{}
	s.Serializer = r.serializer(conn)
	return &s
}

func (r *Rpc) serializer(conn io.ReadWriteCloser) *stream.Serializer {
	s := stream.NewSerializer(conn)
	s.Codec = r.codec()
	s.Logger(r.Log)
	return s
}
