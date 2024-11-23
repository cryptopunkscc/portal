package apphost

import (
	"encoding/json"
	"github.com/cryptopunkscc/astrald/auth/id"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream"
	"io"
)

func RpcClient(
	identity id.Identity,
	query string,
) (s rpc.Conn, err error) {
	return Rpc(Client).Client(identity, query)
}

func (r RpcBase) Client(
	identity id.Identity,
	query string,
) (s rpc.Conn, err error) {
	conn, err := r.client.Query(identity, query)
	if err != nil {
		return
	}
	return NewClient(conn), nil
}

func NewClient(conn io.ReadWriteCloser) *stream.Client {
	s := stream.Client{Serializer: stream.NewSerializer(conn)}
	s.Marshal = json.Marshal
	s.Unmarshal = json.Unmarshal
	return &s
}
