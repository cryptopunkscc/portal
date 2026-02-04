package bind

const (
	ServiceRegister = "_astral_service_register"
	ServiceClose    = "_astral_service_close"
	ConnAccept      = "_astral_conn_accept"
	ConnClose       = "_astral_conn_close"
	ConnWrite       = "_astral_conn_write"
	ConnRead        = "_astral_conn_read"
	ConnWriteLn     = "_astral_conn_write_ln"
	ConnReadLn      = "_astral_conn_read_ln"
	Query           = "_astral_query"
	ResolveId       = "_astral_resolve"
	GetNodeInfo     = "_astral_node_info"
	Interrupt       = "_astral_interrupt"
)

type Apphost interface {
	ServiceRegister() (err error)
	ServiceClose() (err error)
	ConnAccept() (data *QueryData, err error)
	ConnClose(id string) (err error)
	ConnWrite(id string, data []byte) (l int, err error)
	ConnRead(id string, n int) (data []byte, err error)
	ConnWriteLn(id string, data string) (err error)
	ConnReadLn(id string) (data string, err error)
	Query(target string, query string) (data *QueryData, err error)
	QueryString(target string, query string) (data string, err error)
	Resolve(name string) (id string, err error)
	NodeInfo(id string) (info *NodeInfo, err error)
	NodeInfoString(id string) (info string, err error)
	Close() error
	Interrupt()
}

type QueryData struct {
	Id       string `json:"id"`
	Query    string `json:"query"`
	RemoteId string `json:"remoteId"`
}

type NodeInfo struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
}
