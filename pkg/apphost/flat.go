package apphost

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/google/uuid"
	"io"
	"log"
	"sync"
	"time"
)

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
	Resolve         = "_astral_resolve"
)

type FlatAdapter struct {
	listeners      map[string]*astral.Listener
	listenersMutex sync.RWMutex

	connections      map[string]io.ReadWriteCloser
	connectionsMutex sync.RWMutex
}

func NewFlatAdapter() *FlatAdapter {
	return &FlatAdapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]io.ReadWriteCloser{},
	}
}

func (api *FlatAdapter) Interrupt() {
	api.listenersMutex.Lock()
	api.connectionsMutex.Lock()
	defer api.listenersMutex.Unlock()
	defer api.connectionsMutex.Unlock()
	for name, closer := range api.listeners {
		log.Println("[Interrupt] closing listener:", name)
		_ = closer.Close()
	}
	for name, closer := range api.connections {
		log.Print("[Interrupt] closing connection:", name)
		_ = closer.Close()
	}
	api.connections = map[string]io.ReadWriteCloser{}
	api.listeners = map[string]*astral.Listener{}
}

func (api *FlatAdapter) getListener(service string) (l *astral.Listener, ok bool) {
	api.listenersMutex.RLock()
	defer api.listenersMutex.RUnlock()
	l, ok = api.listeners[service]
	return
}

func (api *FlatAdapter) setListener(service string, listener *astral.Listener) {
	api.listenersMutex.Lock()
	defer api.listenersMutex.Unlock()
	if listener != nil {
		api.listeners[service] = listener
	} else {
		delete(api.listeners, service)
	}
}

func (api *FlatAdapter) getConnection(connectionId string) (rw io.ReadWriteCloser, ok bool) {
	api.connectionsMutex.RLock()
	defer api.connectionsMutex.RUnlock()
	rw, ok = api.connections[connectionId]
	return
}

func (api *FlatAdapter) setConnection(connectionId string, connection io.ReadWriteCloser) {
	api.connectionsMutex.Lock()
	defer api.connectionsMutex.Unlock()
	if connection != nil {
		api.connections[connectionId] = connection
	} else {
		delete(api.connections, connectionId)
	}
}

func (api *FlatAdapter) Log(arg ...any) {
	log.Println(arg...)
}

func (api *FlatAdapter) LogArr(arg []any) {
	log.Println(arg...)
}

func (api *FlatAdapter) Sleep(duration int64) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

func (api *FlatAdapter) ServiceRegister(service string) (err error) {
	listener, err := astral.Register(service)
	if err != nil {
		return
	}
	api.setListener(service, listener)
	return
}

func (api *FlatAdapter) ServiceClose(service string) (err error) {
	listener, ok := api.getListener(service)
	if !ok {
		err = errors.New("[ServiceClose] not listening on port: " + service)
		return
	}
	err = listener.Close()
	if err != nil {
		api.setListener(service, nil)
	}
	return
}

func (api *FlatAdapter) ConnAccept(service string) (id string, err error) {
	listener, ok := api.getListener(service)
	if !ok {
		err = fmt.Errorf("[ConnAccept] not listening on port: %v", service)
		return
	}
	conn, err := listener.Accept()
	if err != nil {
		return
	}
	id = uuid.New().String()
	api.setConnection(id, conn)
	return
}

func (api *FlatAdapter) ConnClose(id string) (err error) {
	conn, ok := api.getConnection(id)
	if !ok {
		err = errors.New("[ConnClose] not found connection with id: " + id)
		return
	}
	err = conn.Close()
	if err == nil {
		api.setConnection(id, nil)
	}
	return
}

func (api *FlatAdapter) ConnWrite(id string, data string) (err error) {
	conn, ok := api.getConnection(id)
	if !ok {
		err = errors.New("[ConnWrite] not found connection with id: " + id)
		return
	}
	_, err = conn.Write([]byte(data))
	return
}

func (api *FlatAdapter) ConnRead(id string) (data string, err error) {
	conn, ok := api.getConnection(id)
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

func (api *FlatAdapter) Query(identity string, query string) (connId string, err error) {
	nid := id.Identity{}
	if len(identity) > 0 {
		nid, err = id.ParsePublicKeyHex(identity)
		if err != nil {
			return
		}
	}
	conn, err := astral.Query(nid, query)
	if err != nil {
		return
	}
	connId = uuid.New().String()
	api.setConnection(connId, conn)
	return
}

func (api *FlatAdapter) QueryName(name string, query string) (connId string, err error) {
	conn, err := astral.QueryName(name, query)
	if err != nil {
		return
	}
	connId = uuid.New().String()
	api.setConnection(connId, conn)
	return
}

func (api *FlatAdapter) Resolve(name string) (id string, err error) {
	identity, err := astral.Resolve(name)
	if err != nil {
		return
	}
	id = identity.String()
	return
}

func (api *FlatAdapter) NodeInfo(identity string) (info NodeInfo, err error) {
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
