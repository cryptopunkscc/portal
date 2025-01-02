package apphost

import (
	"encoding/json"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream"
	"io"
)

func RpcClient(
	identity id.Identity,
	service string,
	path ...string,
) (s rpc.Conn, err error) {
	conn, err := Client.Query(identity, apphost.FormatPort(service, path...))
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