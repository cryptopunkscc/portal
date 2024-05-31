package target

import "github.com/cryptopunkscc/astrald/sig"

const (
	Log             = "_log"
	Sleep           = "_sleep"
	ServiceRegister = "_astral_service_register"
	ServiceClose    = "_astral_service_close"
	ConnAccept      = "_astral_conn_accept"
	ConnClose       = "_astral_conn_close"
	ConnWrite       = "_astral_conn_write"
	ConnRead        = "_astral_conn_read"
	Query           = "_astral_query"
	QueryName       = "_astral_query_name"
	GetNodeInfo     = "_astral_node_info"
	ResolveId       = "_astral_resolve"
	Interrupt       = "_astral_interrupt"
)

type Apphost interface {
	ApphostCache
	ApphostApi
}

type ApphostApi interface {
	Prefix() []string
	Close() error
	Interrupt()
	Log(arg ...any)
	LogArr(arg []any)
	Sleep(duration int64)
	ServiceRegister(service string) (err error)
	ServiceClose(service string) (err error)
	ConnAccept(service string) (data string, err error)
	ConnClose(id string) (err error)
	ConnWrite(id string, data string) (err error)
	ConnRead(id string) (data string, err error)
	Query(identity string, query string) (data string, err error)
	QueryName(name string, query string) (data string, err error)
	Resolve(name string) (id string, err error)
	NodeInfo(identity string) (info NodeInfo, err error)
}

type ApphostCache interface {
	Connections() []ApphostConn
	Listeners() []ApphostListener
	Events() *sig.Queue[ApphostEvent]
}

type ApphostEvent struct {
	Type ApphostEventType
	Port string
	Ref  string
}

type ApphostEventType int

const (
	ApphostConnect ApphostEventType = iota
	ApphostDisconnect
	ApphostRegister
	ApphostUnregister
)

type NodeInfo struct {
	Identity string
	Name     string
}

type ApphostConn struct {
	Query string
	In    bool
}

type ApphostListener struct {
	Port string
}
