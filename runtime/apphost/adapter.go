package apphost

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/target"
	"github.com/google/uuid"
	"io"
	"strings"
	"time"
)

func NewAdapter(ctx context.Context, pkg string) target.Apphost {
	if pkg == "" {
		panic("package is empty")
	}
	a := &Adapter{}
	a.port = port.New(pkg)
	a.log = plog.Get(ctx).Type(a).Set(&ctx)

	a.listeners = mem.NewCache[*Listener]()
	a.connections = mem.NewCache[*Conn]()

	a.listeners.OnChange(eventEmitter[*Listener](a.Events()))
	a.connections.OnChange(eventEmitter[*Conn](a.Events()))

	return a
}

func eventEmitter[T any](queue *sig.Queue[target.ApphostEvent]) func(ref string, conn T, added bool) {
	return func(ref string, conn T, added bool) {
		event := target.ApphostEvent{Ref: ref}
		switch v := any(conn).(type) {
		case *Conn:
			event.Port = v.conn.Query()
			event.Type = target.ApphostDisconnect
			if added {
				event.Type = target.ApphostConnect
			}
		case *Listener:
			event.Port = v.port
			event.Type = target.ApphostUnregister
			if added {
				event.Type = target.ApphostRegister
			}
		default:
			return
		}
		queue.Push(event)
	}
}

type Adapter struct {
	log plog.Logger

	port port.Port

	listeners   mem.Cache[*Listener]
	connections mem.Cache[*Conn]
	events      sig.Queue[target.ApphostEvent]
}

func (api *Adapter) Connections() (c []target.ApphostConn) {
	for _, s := range api.connections.Copy() {
		c = append(c, target.ApphostConn{Query: s.conn.Query(), In: s.in})
	}
	return
}

func (api *Adapter) Listeners() (l []target.ApphostListener) {
	for _, s := range api.listeners.Copy() {
		l = append(l, target.ApphostListener{Port: s.port})
	}
	return
}

func (api *Adapter) Events() *sig.Queue[target.ApphostEvent] {
	return &api.events
}

func (api *Adapter) Port(service ...string) string {
	return port.New(service...).String()
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
		port = api.port.String() + service
	case "":
		port = api.port.String()
	default:
		port = api.port.Add(service).String()
	}
	api.log.Println("register:", service)
	astralListener, err := astral.Register(port)
	if err != nil {
		return
	}
	listener := newListener(astralListener, port)
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
		api.listeners.Delete(service)
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
	api.connections.Set(connId, newConn(conn, false))

	base := api.port.String()
	query := strings.TrimPrefix(conn.Query(), base)
	query = strings.TrimPrefix(query, ".")
	bytes, err := json.Marshal(queryData{
		Id:       connId,
		Query:    query,
		RemoteId: conn.RemoteIdentity().String(),
	})
	if err != nil {
		return
	}
	data = string(bytes)
	return
}

type queryData struct {
	Id       string `json:"id"`
	Query    string `json:"query"`
	RemoteId string `json:"remoteId"`
}

func (api *Adapter) ConnClose(id string) (err error) {
	conn, ok := api.connections.Get(id)
	if !ok {
		err = errors.New("[ConnClose] not found connection with id: " + id)
		return
	}
	err = conn.Close()
	api.connections.Delete(id)
	return
}

func (api *Adapter) ConnWrite(id string, data string) (err error) {
	api.log.Printf("> %s <%s>", strings.TrimRight(data, "\r\n"), id)
	conn, ok := api.connections.Get(id)
	if !ok {
		err = errors.New("[ConnWrite] not found connection with id: " + id)
		return
	}
	if _, err = conn.Write([]byte(data)); err != nil {
		_ = conn.Close()
		api.connections.Delete(id)
	}
	return
}

func (api *Adapter) ConnRead(id string) (data string, err error) {
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
	api.log.Printf("< %s <%s>", data, id)
	return
}

func (api *Adapter) Query(identity string, query string) (data string, err error) {
	api.log.Println("~>", query)
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
	api.connections.Set(connId, newConn(conn, false))

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
	api.connections.Set(connId, newConn(conn, false))

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
	in bool
}

type Listener struct {
	*astral.Listener
	port string
}

func newListener(listener *astral.Listener, port string) *Listener {
	return &Listener{Listener: listener, port: port}
}

func newConn(conn *astral.Conn, in bool) *Conn {
	return &Conn{
		conn:        conn,
		Reader:      bufio.NewReader(conn),
		WriteCloser: conn,
		in:          in,
	}
}
