package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"io"
	"net"
)

type Client interface {
	Query(target string, method string, args any) (Conn, error)
	Resolve(name string) (*astral.Identity, error)
	Register() (Listener, error)
	Protocol() string
	DisplayName(identity *astral.Identity) string
	Session() (Session, error)
	Rpc() rpc.Rpc
}

type Session interface {
	Token(token string) (res TokenResponse, err error)
	Query(callerID *astral.Identity, targetID *astral.Identity, query string) (conn Conn, err error)
	Register(identity *astral.Identity, target string) (token string, err error)
	Close() error
}

type TokenResponse interface {
	Code() uint8
	GuestID() *astral.Identity
	HostID() *astral.Identity
}

type Conn interface {
	io.ReadWriteCloser
	RemoteIdentity() *astral.Identity
	RemoteAddr() net.Addr
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
	Token() string
	SetToken(token string)
	Done() <-chan struct{}
}

type PendingQuery interface {
	Query() string
	RemoteIdentity() *astral.Identity
	Reject() (err error)
	Accept() (conn Conn, err error)
	Close() error
}
