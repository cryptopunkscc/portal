package bind

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	"strings"
)

func Adapter(ctx context.Context, cached apphost.Cached, pkg string) Apphost {
	if pkg == "" {
		panic("package is empty")
	}
	a := &adapter{}
	a.Cached = cached
	a.port = port.New(pkg)
	a.log = plog.Get(ctx).Type(a).Set(&ctx)
	return a
}

type adapter struct {
	apphost.Cached
	log  plog.Logger
	port port.Port
}

func (api *adapter) Port(service ...string) string {
	return port.New(service...).String()
}

func (api *adapter) Close() error {
	api.Cached.Interrupt()
	return nil
}

func (api *adapter) ServiceRegister(service string) (err error) {
	service = api.localPort(service)
	api.log.Println("register:", service)
	_, err = api.Cached.Register(service)
	return
}

func (api *adapter) ServiceClose(service string) (err error) {
	service = api.localPort(service)
	listener, ok := api.Listeners().Get(service)
	if !ok {
		err = errors.New("[ServiceClose] not listening on port: " + service)
		return
	}
	return listener.Close()
}

func (api *adapter) ConnAccept(service string) (data string, err error) {
	service = api.localPort(service)
	listener, ok := api.Listeners().Get(service)
	if !ok {
		err = fmt.Errorf("[ConnAccept] not listening on port: %v, %v", service, api.Listeners().Copy())
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

	api.log.Println("accepted connection:", conn.Ref())

	base := api.port.String()
	query := strings.TrimPrefix(conn.Query(), base)
	query = strings.TrimPrefix(query, ".")
	bytes, err := json.Marshal(queryData{
		Id:       conn.Ref(),
		Query:    query,
		RemoteId: conn.RemoteIdentity().String(),
	})
	if err != nil {
		return
	}
	data = string(bytes)
	return
}

func (api *adapter) localPort(service string) string {
	switch service {
	case "*":
		return api.port.String() + service
	case "":
		return api.port.String()
	}
	return api.port.Add(service).String()
}

type queryData struct {
	Id       string `json:"id"`
	Query    string `json:"query"`
	RemoteId string `json:"remoteId"`
}

func (api *adapter) ConnClose(id string) (err error) {
	conn, ok := api.Connections().Get(id)
	if !ok {
		err = errors.New("[ConnClose] not found connection with id: " + id)
		return
	}
	err = conn.Close()
	return
}

func (api *adapter) ConnWrite(id string, data string) (err error) {
	api.log.Printf("> %s <%s>", strings.TrimRight(data, "\r\n"), id)
	conn, ok := api.Connections().Get(id)
	if !ok {
		err = errors.New("[ConnWrite] not found connection with id: " + id)
		return
	}
	if _, err = conn.Write([]byte(data)); err != nil {
		_ = conn.Close()
	}
	return
}

func (api *adapter) ConnRead(id string) (data string, err error) {
	conn, ok := api.Connections().Get(id)
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

func (api *adapter) Query(identity string, query string) (data string, err error) {
	api.log.Println("~>", query)
	nid := id.Identity{}
	if len(identity) > 0 {
		nid, err = id.ParsePublicKeyHex(identity)
		if err != nil {
			return
		}
	}
	conn, err := api.Cached.Query(nid, api.Port(query))
	if err != nil {
		return
	}

	bytes, err := json.Marshal(queryData{
		Id:    conn.Ref(),
		Query: conn.Query(),
	})
	if err != nil {
		return
	}
	data = string(bytes)
	return
}

func (api *adapter) QueryName(name string, query string) (data string, err error) {
	conn, err := api.Cached.QueryName(name, api.Port(query))
	if err != nil {
		return
	}

	bytes, err := json.Marshal(queryData{
		Id:    conn.Ref(),
		Query: conn.Query(),
	})
	if err != nil {
		return
	}
	data = string(bytes)
	return
}

func (api *adapter) Resolve(name string) (id string, err error) {
	identity, err := api.Cached.Resolve(name)
	if err != nil {
		return
	}
	id = identity.String()
	return
}

func (api *adapter) NodeInfo(identity string) (info *bind.NodeInfo, err error) {
	nid, err := id.ParsePublicKeyHex(identity)
	if err != nil {
		return
	}
	i, err := api.Cached.NodeInfo(nid)
	if err != nil {
		return
	}
	info = &bind.NodeInfo{
		Identity: i.Identity.String(),
		Name:     i.Name,
	}
	return
}
