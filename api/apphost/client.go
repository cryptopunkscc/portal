package apphost

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/astrald/mod/apphost/proto"
	"io"
	"net"
)

type Client interface {
	Discovery() Discovery
	Session() (Session, error)
	Query(remoteID id.Identity, query string) (conn Conn, err error)
	QueryName(name string, query string) (conn Conn, err error)
	Resolve(name string) (id.Identity, error)
	NodeInfo(identity id.Identity) (info proto.NodeInfoData, err error)
	Exec(identity id.Identity, app string, args []string, env []string) error
	Register(service string) (l Listener, err error)
}

type Discovery interface {
	Discover(identity id.Identity) ([]astral.ServiceInfo, error)
}

type Session interface {
	Query(remoteID id.Identity, query string) (conn Conn, err error)
	Resolve(name string) (identity id.Identity, err error)
	NodeInfo(identity id.Identity) (info proto.NodeInfoData, err error)
	Register(service string, target string) (err error)
	Exec(identity id.Identity, app string, args []string, env []string) (err error)
}

type Conn interface {
	io.ReadWriteCloser
	RemoteIdentity() id.Identity
	RemoteAddr() net.Addr
	Query() string
	ReadString(delim byte) (string, error)

	Ref() string
	In() bool
}

type Listener interface {
	Next() (QueryData, error)
	QueryCh() <-chan QueryData
	Accept() (net.Conn, error)
	AcceptAll() <-chan net.Conn
	Close() error
	Addr() net.Addr
	Target() string

	Port() string
}

type QueryData interface {
	Query() string
	RemoteIdentity() id.Identity
	Reject() error
	Accept() (Conn, error)
}
