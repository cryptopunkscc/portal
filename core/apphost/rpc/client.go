package rpc

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream/query"
	"io"
)

func (r Rpc) Client(
	target string,
	query string,
) (s rpc.Conn, err error) {
	conn, err := r.Apphost.Query(target, query, nil)
	if err != nil {
		return
	}
	return rpcClient(conn), nil
}

func rpcClient(conn io.ReadWriteCloser) *stream.Client {
	s := stream.Client{Serializer: stream.NewSerializer(conn)}
	s.MarshalArgs = query.Marshal
	s.Marshal = json.Marshal
	s.Unmarshal = json.Unmarshal
	s.Ending = []byte("\n")
	return &s
}
