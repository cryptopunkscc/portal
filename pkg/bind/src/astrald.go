package bind

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	libquery "github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/portal/pkg/bind"
	"github.com/cryptopunkscc/portal/pkg/client"
)

type Astrald struct {
	client.Astrald
	cache
	listener *astrald.Listener
}

func (a *Astrald) Interrupt() {
	_ = a.Close()
}

func (a *Astrald) Close() error {
	_ = a.ServiceClose()
	a.interrupt()
	return nil
}

func (a *Astrald) ServiceRegister() (err error) {
	if a.listener != nil {
		return errors.New("already listening")
	}
	a.listener, err = a.Register(context.Background())
	return
}

func (a *Astrald) ServiceClose() (err error) {
	listener := a.listener
	if listener == nil {
		return
	}
	err = listener.Close()
	a.listener = nil
	return
}

func (a *Astrald) ConnAccept() (data *bind.QueryData, err error) {
	listener := a.listener
	if listener == nil {
		return nil, fmt.Errorf("[ConnAccept] not listening: %v", listener)
	}

	var next *astrald.PendingQuery
	if next, err = listener.Next(); err != nil {
		return
	}

	conn := next.Accept()
	a.set(conn.Query(), conn)

	data = &bind.QueryData{
		Id:       conn.Query().Nonce.String(),
		Query:    conn.Query().Query,
		RemoteId: conn.RemoteIdentity().String(),
	}

	a.Log.Println("accepted connection:", data.Id)
	return
}

func (a *Astrald) ConnClose(id string) (err error) {
	a.Log.Printf("close <%s>", id)
	c, ok := a.get(id)
	if !ok {
		return errors.New("[ConnClose] not found connection with id: " + id)
	}
	return c.Close()
}

func (a *Astrald) ConnWrite(id string, data []byte) (n int, err error) {
	a.Log.Printf("> [%v]byte <%s>", len(data), id)
	//api.Log.Printf("> [%v]byte <%s>", data, id)

	conn, ok := a.get(id)
	if !ok {
		err = errors.New("[ConnWrite] not found connection with id: " + id)
		return
	}
	n, err = conn.Write(data)
	return
}

func (a *Astrald) ConnRead(id string, n int) (data []byte, err error) {
	conn, ok := a.get(id)
	if !ok {
		err = errors.New("[ConnRead] not found connection with id: " + id)
		return
	}
	buf := make([]byte, n)
	n, err = conn.Read(buf)
	data = buf[:n]
	a.Log.Printf("< [%v]byte <%s>", n, id)
	//api.Log.Printf("< [%v]byte <%s> %v", data, id, err)
	return
}

func (a *Astrald) ConnWriteLn(id string, data string) (err error) {
	a.Log.Printf("> %s <%s>", strings.TrimRight(data, "\r\n"), id)
	conn, ok := a.get(id)
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

func (a *Astrald) ConnReadLn(id string) (data string, err error) {
	conn, ok := a.get(id)
	if !ok {
		err = errors.New("[ConnReadLn] not found connection with id: " + id)
		return
	}
	data, err = conn.ReadString('\n')
	data = strings.TrimSuffix(data, "\n")
	a.Log.Printf("< %s <%s>", data, id)
	return
}

func (a *Astrald) Query(target string, query string) (data *bind.QueryData, err error) {
	a.Log.Println("~>", target, query)
	targetID, err := a.Astrald.Resolve(target)
	if err != nil {
		return
	}
	q := libquery.New(a.GuestID(), targetID, query, nil)
	conn, err := a.RouteQuery(astral.NewContext(nil), q)
	if err != nil {
		return
	}
	a.set(q, conn)
	data = &bind.QueryData{
		Id:       q.Nonce.String(),
		Query:    query,
		RemoteId: conn.RemoteIdentity().String(),
	}
	return
}

func (a *Astrald) QueryString(target string, query string) (data string, err error) {
	queryData, err := a.Query(target, query)
	if err != nil {
		return
	}
	bytes, err := json.Marshal(queryData)
	if err != nil {
		return
	}
	data = string(bytes)
	return
}

func (a *Astrald) Resolve(name string) (id string, err error) {
	identity, err := a.Astrald.Resolve(name)
	if err != nil {
		return
	}
	id = identity.String()
	return
}

func (a *Astrald) NodeInfo(identity string) (info *bind.NodeInfo, err error) {
	nid, err := astral.ParseIdentity(identity)
	if err != nil {
		return
	}
	alias, err := a.Dir().GetAlias(nil, nid)
	if err != nil {
		return
	}
	info = &bind.NodeInfo{
		Identity: identity,
		Name:     alias,
	}
	return
}

func (a *Astrald) NodeInfoString(identity string) (info string, err error) {
	nodeInfo, err := a.NodeInfo(identity)
	if err != nil {
		return
	}
	bytes, err := json.Marshal(nodeInfo)
	if err != nil {
		return
	}
	info = string(bytes)
	return
}
