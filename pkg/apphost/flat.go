package apphost

import (
	"bufio"
	_ "embed"
	"encoding/json"
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
	Interrupt       = "_astral_interrupt"
)

type FlatAdapter struct {
	listeners      map[string]*astral.Listener
	listenersMutex sync.RWMutex

	connections      map[string]*Conn
	connectionsMutex sync.RWMutex
}

func NewFlatAdapter() *FlatAdapter {
	return &FlatAdapter{
		listeners:   map[string]*astral.Listener{},
		connections: map[string]*Conn{},
	}
}

func (api *FlatAdapter) Close() error {
	api.Interrupt()
	return nil
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
	api.connections = map[string]*Conn{}
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

func (api *FlatAdapter) getConnection(connectionId string) (rw *Conn, ok bool) {
	api.connectionsMutex.RLock()
	defer api.connectionsMutex.RUnlock()
	rw, ok = api.connections[connectionId]
	return
}

func (api *FlatAdapter) setConnection(connectionId string, connection *astral.Conn) {
	api.connectionsMutex.Lock()
	defer api.connectionsMutex.Unlock()
	if connection != nil {
		api.connections[connectionId] = newConn(connection)
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

func (api *FlatAdapter) ConnAccept(service string) (data string, err error) {
	listener, ok := api.getListener(service)
	if !ok {
		err = fmt.Errorf("[ConnAccept] not listening on port: %v", service)
		return
	}
	next, err := listener.Next()
	if err != nil {
		return
	}

	conn, err := next.Accept()
	if err != nil {
		return
	}

	connId := uuid.New().String()
	api.setConnection(connId, conn)

	bytes, err := json.Marshal(queryData{
		Id:    connId,
		Query: conn.Query(),
	})
	if err != nil {
		return
	}
	data = string(bytes)
	return
}

type queryData struct {
	Id    string `json:"id"`
	Query string `json:"query"`
}

func (api *FlatAdapter) ConnClose(id string) (err error) {
	conn, ok := api.getConnection(id)
	if !ok {
		err = errors.New("[ConnClose] not found connection with id: " + id)
		return
	}
	err = conn.Close()
	api.setConnection(id, nil)
	return
}

func (api *FlatAdapter) ConnWrite(id string, data string) (err error) {
	conn, ok := api.getConnection(id)
	if !ok {
		err = errors.New("[ConnWrite] not found connection with id: " + id)
		return
	}
	if _, err = conn.Write([]byte(data)); err != nil {
		_ = conn.Close()
		api.setConnection(id, nil)
	}
	return
}

func (api *FlatAdapter) ConnRead(id string) (data string, err error) {
	conn, ok := api.getConnection(id)
	if !ok {
		err = errors.New("[ConnRead] not found connection with id: " + id)
		return
	}
	return conn.ReadString('\n')
}

func (api *FlatAdapter) Query(identity string, query string) (data string, err error) {
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
	connId := uuid.New().String()
	api.setConnection(connId, conn)

	bytes, err := json.Marshal(queryData{
		Id:    connId,
		Query: conn.Query(),
	})
	if err != nil {
		return
	}
	data = string(bytes)
	return
}

func (api *FlatAdapter) QueryName(name string, query string) (data string, err error) {
	conn, err := astral.QueryName(name, query)
	if err != nil {
		return
	}
	connId := uuid.New().String()
	api.setConnection(connId, conn)

	bytes, err := json.Marshal(queryData{
		Id:    connId,
		Query: conn.Query(),
	})
	if err != nil {
		return
	}
	data = string(bytes)
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

type Conn struct {
	conn *astral.Conn
	*bufio.Reader
	io.WriteCloser
}

func newConn(conn *astral.Conn) *Conn {
	return &Conn{
		conn:        conn,
		Reader:      bufio.NewReader(conn),
		WriteCloser: conn,
	}
}

type NodeInfo struct {
	Identity string
	Name     string
}
