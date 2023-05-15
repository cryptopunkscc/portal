package astral_js

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

//go:embed apphost.js
var appHostJsClient string

func AppHostJsClient() string {
	return appHostJsClient
}

type AppHostFlatAdapter struct {
	ports map[string]*astral.Listener
	conns map[string]io.ReadWriteCloser
}

func NewAppHostFlatAdapter() *AppHostFlatAdapter {
	return &AppHostFlatAdapter{
		ports: map[string]*astral.Listener{},
		conns: map[string]io.ReadWriteCloser{},
	}
}

func (api *AppHostFlatAdapter) Log(arg ...any) {
	log.Println(arg...)
}

func (api *AppHostFlatAdapter) Sleep(duration int64) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

func (api *AppHostFlatAdapter) PortListen(port string) (err error) {
	listener, err := astral.Listen(port)
	if err != nil {
		return
	}
	api.ports[port] = listener
	return
}

func (api *AppHostFlatAdapter) PortClose(port string) (err error) {
	listener, ok := api.ports[port]
	if !ok {
		err = errors.New("[PortClose] not listening on port: " + port)
		return
	}
	err = listener.Close()
	if err != nil {
		delete(api.ports, port)
	}
	return
}

func (api *AppHostFlatAdapter) ConnAccept(port string) (id string, err error) {
	listener, ok := api.ports[port]
	if !ok {
		err = fmt.Errorf("[ConnAccept] not listening on port: %v", port)
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

func (api *AppHostFlatAdapter) Dial(identity string, query string) (connId string, err error) {
	nid, err := id.ParsePublicKeyHex(identity)
	if err != nil {
		return
	}
	conn, err := astral.Dial(nid, query)
	if err != nil {
		return
	}
	connId = uuid.New().String()
	api.conns[connId] = conn
	return
}

func (api *AppHostFlatAdapter) DialName(name string, query string) (connId string, err error) {
	conn, err := astral.DialName(name, query)
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
	i, err := astral.NodeInfo(nid)
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
