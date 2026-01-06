package apphost

import (
	"io"
	"net"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

type Client interface {
	Query(target string, method string, args any) (Conn, error)
	Resolve(name string) (*astral.Identity, error)
	Register() (Listener, error)
	DisplayName(identity *astral.Identity) string
	Rpc() rpc.Rpc
}

type Conn interface {
	io.ReadWriteCloser
	LocalIdentity() *astral.Identity
	RemoteIdentity() *astral.Identity
	Query() string
	ReadString(delim byte) (string, error)

	Ref() string
	In() bool
}

type Listener interface {
	Next() (PendingQuery, error)
	Accept() (net.Conn, error)
	Close() error
	Addr() net.Addr
	String() string
	Done() <-chan struct{}
}

type PendingQuery interface {
	Query() string
	Caller() *astral.Identity
	Skip() error
	Reject() (err error)
	Accept() (conn Conn, err error)
	Close() error
}
