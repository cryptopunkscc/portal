package apphost

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/runtime/rpc"
	"github.com/cryptopunkscc/portal/runtime/rpc/stream"
	"github.com/cryptopunkscc/portal/runtime/rpc/stream/query"
	"io"
)

func (r RpcBase) Client(
	target string,
	query string,
) (s rpc.Conn, err error) {
	conn, err := r.client.Query(target, query, nil)
	if err != nil {
		return
	}
	return NewClient(conn), nil
}

func NewClient(conn io.ReadWriteCloser) *stream.Client {
	s := stream.Client{Serializer: stream.NewSerializer(conn)}
	s.MarshalArgs = query.Marshal
	s.Marshal = json.Marshal
	s.Unmarshal = json.Unmarshal
	s.Ending = []byte("\n")
	return &s
}
