package bind

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"strings"
)

func Adapter(ctx context.Context, cached apphost.Cached, pkg string) Apphost {
	if pkg == "" {
		panic("package is empty")
	}
	a := &adapter{}
	a.Cached = cached
	a.pkg = pkg
	a.log = plog.Get(ctx).Type(a)
	return a
}

type adapter struct {
	apphost.Cached
	log      plog.Logger
	pkg      string
	listener apphost.Listener
}

func (a *adapter) Close() error {
	_ = a.ServiceClose()
	a.Cached.Interrupt()
	return nil
}

func (a *adapter) ServiceRegister() (err error) {
	a.listener, err = a.Cached.Register()
	return
}

func (a *adapter) ServiceClose() (err error) {
	listener := a.listener
	if listener == nil {
		return
	}
	err = listener.Close()
	a.listener = nil
	return
}

func (a *adapter) ConnAccept() (data string, err error) {
	listener := a.listener
	if listener == nil {
		err = fmt.Errorf("[ConnAccept] not listening: %v", listener)
		return
	}

	var next apphost.PendingQuery
	for {
		if next, err = listener.Next(); err != nil {
			return
		}
		if strings.HasPrefix(next.Query(), a.pkg) {
			break
		}
		_ = next.Close()
	}
	conn, err := next.Accept()
	if err != nil {
		return
	}

	a.log.Println("accepted connection:", conn.Ref())

	bytes, err := json.Marshal(queryData{
		Id:       conn.Ref(),
		Query:    conn.Query(),
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

func (a *adapter) ConnClose(id string) (err error) {
	a.log.Printf("close <%s>", id)
	conn, ok := a.Connections().Get(id)
	if !ok {
		err = errors.New("[ConnClose] not found connection with id: " + id)
		return
	}
	err = conn.Close()
	return
}

func (a *adapter) ConnWrite(id string, data []byte) (n int, err error) {
	a.log.Printf("> [%v]byte <%s>", len(data), id)
	//api.log.Printf("> [%v]byte <%s>", data, id)
	conn, ok := a.Connections().Get(id)
	if !ok {
		err = errors.New("[ConnWrite] not found connection with id: " + id)
		return
	}
	n, err = conn.Write(data)
	return
}

func (a *adapter) ConnRead(id string, n int) (data []byte, err error) {
	conn, ok := a.Connections().Get(id)
	if !ok {
		err = errors.New("[ConnRead] not found connection with id: " + id)
		return
	}
	buf := make([]byte, n)
	n, err = conn.Read(buf)
	data = buf[:n]
	a.log.Printf("< [%v]byte <%s>", n, id)
	//api.log.Printf("< [%v]byte <%s> %v", data, id, err)
	return
}

func (a *adapter) ConnWriteLn(id string, data string) (err error) {
	a.log.Printf("> %s <%s>", strings.TrimRight(data, "\r\n"), id)
	conn, ok := a.Connections().Get(id)
	if !ok {
		err = errors.New("[ConnWriteLn] not found connection with id: " + id)
		return
	}
	if !strings.HasSuffix(data, "\n") {
		data += "\n"
	}
	_, err = conn.Write([]byte(data))
	return
}

func (a *adapter) ConnReadLn(id string) (data string, err error) {
	conn, ok := a.Connections().Get(id)
	if !ok {
		err = errors.New("[ConnReadLn] not found connection with id: " + id)
		return
	}
	data, err = conn.ReadString('\n')
	data = strings.TrimSuffix(data, "\n")
	a.log.Printf("< %s <%s>", data, id)
	return
}

func (a *adapter) Query(target string, query string) (data string, err error) {
	a.log.Println("~>", query)
	conn, err := a.Cached.Query(target, query, nil)
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

func (a *adapter) QueryName(name string, query string) (data string, err error) {
	id, err := a.Resolve(name)
	if err != nil {
		return
	}
	return a.Query(id, query)
}

func (a *adapter) Resolve(name string) (id string, err error) {
	identity, err := a.Cached.Resolve(name)
	if err != nil {
		return
	}
	id = identity.String()
	return
}

func (a *adapter) NodeInfo(identity string) (info *bind.NodeInfo, err error) {
	nid, err := astral.IdentityFromString(identity)
	if err != nil {
		return
	}
	name := a.Cached.DisplayName(nid)
	info = &bind.NodeInfo{
		Identity: identity,
		Name:     name,
	}
	return
}
