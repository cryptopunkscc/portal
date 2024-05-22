package apphost

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

type Adapter struct {
	pkg    []string
	prefix []string

	listeners      map[string]*astral.Listener
	listenersMutex sync.RWMutex

	connections      map[string]*Conn
	connectionsMutex sync.RWMutex

	onIdle func(bool)
}

func (api *Adapter) Port(service ...string) (port string) {
	return strings.Join(append(api.Prefix(), service...), ".")
}

func (api *Adapter) Prefix() []string {
	return api.prefix
}

func (api *Adapter) Close() error {
	api.Interrupt()
	return nil
}

func (api *Adapter) Interrupt() {
	api.listenersMutex.Lock()
	api.connectionsMutex.Lock()
	defer api.listenersMutex.Unlock()
	defer api.connectionsMutex.Unlock()
	for name, closer := range api.listeners {
		log.Println("[Interrupt] closing listener:", name)
		_ = closer.Close()
		delete(api.listeners, name)
	}
	for name := range api.connections {
		log.Print("[Interrupt] closing connection:", name)
		api.setConnectionUnsafe(name, nil)
	}
	api.connections = map[string]*Conn{}
	api.listeners = map[string]*astral.Listener{}
}

func (api *Adapter) getListener(service string) (l *astral.Listener, ok bool) {
	api.listenersMutex.RLock()
	defer api.listenersMutex.RUnlock()
	l, ok = api.listeners[service]
	return
}

func (api *Adapter) setListener(service string, listener *astral.Listener) {
	api.listenersMutex.Lock()
	defer api.listenersMutex.Unlock()
	if listener != nil {
		api.listeners[service] = listener
	} else {
		delete(api.listeners, service)
	}
}

func (api *Adapter) getConnection(connectionId string) (rw *Conn, ok bool) {
	api.connectionsMutex.RLock()
	defer api.connectionsMutex.RUnlock()
	rw, ok = api.connections[connectionId]
	return
}

func (api *Adapter) setConnection(connectionId string, connection *astral.Conn) {
	api.connectionsMutex.Lock()
	defer api.connectionsMutex.Unlock()
	api.setConnectionUnsafe(connectionId, connection)
}

func (api *Adapter) setConnectionUnsafe(connectionId string, connection *astral.Conn) {
	if connection != nil {
		api.connections[connectionId] = newConn(connection)
	} else {
		delete(api.connections, connectionId)
	}
	if api.onIdle != nil {
		api.onIdle(len(api.connections) == 0)
	}
}

func (api *Adapter) Log(arg ...any) {
	log.Println(arg...)
}

func (api *Adapter) LogArr(arg []any) {
	log.Println(arg...)
}

func (api *Adapter) Sleep(duration int64) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

func (api *Adapter) ServiceRegister(service string) (err error) {
	port := service
	switch service {
	case "*":
		port = api.Port(api.pkg...) + service
	case "":
		port = api.Port(api.pkg...)
	default:
		port = api.Port(append(api.pkg, service)...)
	}
	listener, err := astral.Register(port)
	if err != nil {
		return
	}
	api.setListener(service, listener)
	return
}

func (api *Adapter) ServiceClose(service string) (err error) {
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

func (api *Adapter) ConnAccept(service string) (data string, err error) {
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

	pkg := strings.Join(append(api.Prefix(), api.pkg...), ".")
	query := strings.TrimPrefix(conn.Query(), pkg)
	query = strings.TrimPrefix(query, ".")
	bytes, err := json.Marshal(queryData{
		Id:    connId,
		Query: query,
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

func (api *Adapter) ConnClose(id string) (err error) {
	conn, ok := api.getConnection(id)
	if !ok {
		err = errors.New("[ConnClose] not found connection with id: " + id)
		return
	}
	err = conn.Close()
	api.setConnection(id, nil)
	return
}

func (api *Adapter) ConnWrite(id string, data string) (err error) {
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

func (api *Adapter) ConnRead(id string) (data string, err error) {
	conn, ok := api.getConnection(id)
	if !ok {
		err = errors.New("[ConnRead] not found connection with id: " + id)
		return
	}
	data, err = conn.ReadString('\n')
	if err != nil {
		return
	}
	data = strings.TrimSuffix(data, "\n")
	return
}

func (api *Adapter) Query(identity string, query string) (data string, err error) {
	nid := id.Identity{}
	if len(identity) > 0 {
		nid, err = id.ParsePublicKeyHex(identity)
		if err != nil {
			return
		}
	}
	conn, err := astral.Query(nid, api.Port(query))
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

func (api *Adapter) QueryName(name string, query string) (data string, err error) {
	conn, err := astral.QueryName(name, api.Port(query))
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

func (api *Adapter) Resolve(name string) (id string, err error) {
	identity, err := astral.Resolve(name)
	if err != nil {
		return
	}
	id = identity.String()
	return
}

func (api *Adapter) NodeInfo(identity string) (info target.NodeInfo, err error) {
	nid, err := id.ParsePublicKeyHex(identity)
	if err != nil {
		return
	}
	i, err := astral.GetNodeInfo(nid)
	if err != nil {
		return
	}
	info = target.NodeInfo{
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
