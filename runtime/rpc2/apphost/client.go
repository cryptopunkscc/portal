package apphost

import (
	"encoding/json"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/port"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"io"
)

func RpcClient(
	identity id.Identity,
	service string,
	path ...string,
) (s rpc.Conn, err error) {
	conn, err := Client.Query(identity, port.Format(service, path...))
	if err != nil {
		return
	}
	return NewClient(conn), nil
}

func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	s := rpc.Client{Serializer: rpc.NewSerializer(conn)}
	s.Marshal = json.Marshal
	s.Unmarshal = json.Unmarshal
	return &s
}
