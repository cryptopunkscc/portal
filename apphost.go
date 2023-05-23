package astraljs

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/google/uuid"
	"io"
	"log"
	"time"
)

const (
	Log             = "log"
	Sleep           = "sleep"
	ServiceRegister = "astral_service_register"
	ServiceClose    = "astral_service_close"
	ConnAccept      = "astral_conn_accept"
	ConnClose       = "astral_conn_close"
	ConnWrite       = "astral_conn_write"
	ConnRead        = "astral_conn_read"
	Query           = "astral_query"
	QueryName       = "astral_query_name"
	GetNodeInfo     = "astral_node_info"
	Resolve         = "astral_resolve"
)

//go:embed apphost.js
var appHostJsClient string

func AppHostJsClient() string {
	return appHostJsClient
}

type AppHostFlatAdapter struct {
	listeners map[string]*astral.Listener
	conns     map[string]io.ReadWriteCloser
}

func NewAppHostFlatAdapter() *AppHostFlatAdapter {
	return &AppHostFlatAdapter{
		listeners: map[string]*astral.Listener{},
		conns:     map[string]io.ReadWriteCloser{},
	}
}

func (api *AppHostFlatAdapter) Log(arg ...any) {
	log.Println(arg...)
}

func (api *AppHostFlatAdapter) LogArr(arg []any) {
	log.Println(arg...)
}

func (api *AppHostFlatAdapter) Sleep(duration int64) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

func (api *AppHostFlatAdapter) ServiceRegister(service string) (err error) {
	listener, err := astral.Register(service)
	if err != nil {
		return
	}
	api.listeners[service] = listener
	return
}

func (api *AppHostFlatAdapter) ServiceClose(service string) (err error) {
	listener, ok := api.listeners[service]
	if !ok {
		err = errors.New("[ServiceClose] not listening on port: " + service)
		return
	}
	err = listener.Close()
	if err != nil {
		delete(api.listeners, service)
	}
	return
}

func (api *AppHostFlatAdapter) ConnAccept(service string) (id string, err error) {
	listener, ok := api.listeners[service]
	if !ok {
		err = fmt.Errorf("[ConnAccept] not listening on port: %v", service)
		return
	}
	conn, err := listener.Accept()
	if err != nil {
		return
	}
	id = uuid.New().String()
	api.conns[id] = conn
	return
}

func (api *AppHostFlatAdapter) ConnClose(id string) (err error) {
	conn, ok := api.conns[id]
	if !ok {
		err = errors.New("[ConnClose] not found connection with id: " + id)
		return
	}
	err = conn.Close()
	if err == nil {
		delete(api.conns, id)
	}
	return
}

func (api *AppHostFlatAdapter) ConnWrite(id string, data string) (err error) {
	conn, ok := api.conns[id]
	if !ok {
		err = errors.New("[ConnWrite] not found connection with id: " + id)
		return
	}
	_, err = conn.Write([]byte(data))
	return
}

func (api *AppHostFlatAdapter) ConnRead(id string) (data string, err error) {
	conn, ok := api.conns[id]
	if !ok {
		err = errors.New("[ConnRead] not found connection with id: " + id)
		return
	}
	buf := make([]byte, 4096)
	arr := make([]byte, 0)
	n := 0
	defer func() {
		data = string(arr)
	}()
	for {
		n, err = conn.Read(buf)
		if err != nil {
			return
		}
		arr = append(arr, buf[0:n]...)
		if n < len(buf) {
			return
		}
	}
}

func (api *AppHostFlatAdapter) Query(identity string, query string) (connId string, err error) {
	nid, err := id.ParsePublicKeyHex(identity)
	if err != nil {
		return
	}
	conn, err := astral.Query(nid, query)
	if err != nil {
		return
	}
	connId = uuid.New().String()
	api.conns[connId] = conn
	return
}

func (api *AppHostFlatAdapter) QueryName(name string, query string) (connId string, err error) {
	conn, err := astral.QueryName(name, query)
	if err != nil {
		return
	}
	connId = uuid.New().String()
	api.conns[connId] = conn
	return
}

func (api *AppHostFlatAdapter) Resolve(name string) (id string, err error) {
	identity, err := astral.Resolve(name)
	if err != nil {
		return
	}
	id = identity.String()
	return
}

func (api *AppHostFlatAdapter) NodeInfo(identity string) (info NodeInfo, err error) {
	nid, err := id.ParsePublicKeyHex(identity)
	if err != nil {
		return
	}
	i, err := astral.GetNodeInfo(nid)
	if err != nil {
		return
	}
	info = NodeInfo{
		Identity: i.Identity.String(),
		Name:     i.Name,
	}
	return
}

type NodeInfo struct {
	Identity string
	Name     string
}
