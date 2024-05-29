package apphost

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/google/uuid"
	"io"
	"strings"
	"time"
)

type Adapter struct {
	log plog.Logger

	pkg    []string
	prefix []string

	listeners   Registry[*astral.Listener]
	connections Registry[*Conn]
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
	for name, closer := range api.listeners.Release() {
		api.log.Println("[Interrupt] closing listener:", name)
		_ = closer.Close()
	}
	for name, conn := range api.connections.Release() {
		api.log.Println("[Interrupt] closing connection:", name)
		_ = conn.Close()
	}
}

func (api *Adapter) Log(arg ...any) {
	api.log.Scope("js").Println(arg...)
}

func (api *Adapter) LogArr(arg []any) {
	api.log.Scope("js").Println(arg...)
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
	api.log.Println("register:", service)
	listener, err := astral.Register(port)
	if err != nil {
		return
	}
	api.listeners.Set(service, listener)
	return
}

func (api *Adapter) ServiceClose(service string) (err error) {
	listener, ok := api.listeners.Get(service)
	if !ok {
		err = errors.New("[ServiceClose] not listening on port: " + service)
		return
	}
	err = listener.Close()
	if err != nil {
		api.listeners.Set(service, nil)
	}
	return
}

func (api *Adapter) ConnAccept(service string) (data string, err error) {
	listener, ok := api.listeners.Get(service)
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
	api.log.Println("accepted connection:", connId)
	api.connections.Set(connId, newConn(conn))

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
	conn, ok := api.connections.Get(id)
	if !ok {
		err = errors.New("[ConnClose] not found connection with id: " + id)
		return
	}
	err = conn.Close()
	api.connections.Set(id, nil)
	return
}

func (api *Adapter) ConnWrite(id string, data string) (err error) {
	api.log.Printf("[%s] write: %s", id, strings.TrimRight(data, "\r\n"))
	conn, ok := api.connections.Get(id)
	if !ok {
		err = errors.New("[ConnWrite] not found connection with id: " + id)
		return
	}
	if _, err = conn.Write([]byte(data)); err != nil {
		_ = conn.Close()
		api.connections.Set(id, nil)
	}
	return
}

func (api *Adapter) ConnRead(id string) (data string, err error) {
	api.log.Printf("[%s] read: %s", id, data)
	conn, ok := api.connections.Get(id)
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
	api.log.Printf("[%s] query: %s", identity, query)
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
	api.connections.Set(connId, newConn(conn))

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
	api.connections.Set(connId, newConn(conn))

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
