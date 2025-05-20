package objects

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/objects"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"io"
)

func Client(rpc rpc.Rpc) Conn {
	return Conn{rpc.Format("json").Request("localnode", "objects")}
}

type Conn struct{ rpc.Conn }

type pushArgs struct{ Size int }

func (c Conn) Push(obj astral.Object) (ok bool, err error) {
	defer plog.TraceErr(&err)
	buf := bytes.NewBuffer(nil)
	if _, err = astral.WriteCanonical(buf, obj); err != nil {
		return
	}
	conn := c.Copy()
	if err = rpc.Call(conn, "push", pushArgs{Size: buf.Len()}); err != nil {
		return
	}
	defer conn.Close()
	if _, err = buf.WriteTo(conn); err != nil {
		return
	}
	err = binary.Read(conn, binary.BigEndian, &ok)
	return
}

type ReadArgs struct {
	ID     astral.ObjectID
	Offset astral.Uint64
	Zone   astral.Zone
}

func (c Conn) Read(args ReadArgs) (r io.ReadCloser, err error) {
	return rpc.NewCall(c.Conn, "read", args)
}

func (c Conn) Fetch(args ReadArgs, obj astral.Object) (err error) {
	b, err := c.Read(args)
	if err != nil {
		return
	}
	t, r, err := astral.OpenCanonical(b)
	if err != nil {
		return err
	}
	if t != obj.ObjectType() {
		return fmt.Errorf("expected object type %s, got %s", obj.ObjectType(), t)
	}
	_, err = obj.ReadFrom(r)
	return
}

type SearchArgs struct {
	Query string `query:"q"`
	Zone  astral.Zone
	Out   string
	Ext   string
}

func (c Conn) Search(args SearchArgs) (out <-chan ObjectResponse[objects.SearchResult], err error) {
	args.Out = "json"
	return rpc.Subscribe[ObjectResponse[objects.SearchResult]](c.Conn, "search", args)
}

type ScanArgs struct {
	Type   string
	Repo   string
	Out    string
	Follow bool
	Zone   astral.Zone
}

func (c Conn) Scan(args ScanArgs) (out <-chan ObjectResponse[objects.SearchResult], err error) {
	args.Out = "json"
	return rpc.Subscribe[ObjectResponse[objects.SearchResult]](c.Conn, "scan", args)
}

type DescribeArgs struct {
	ID    astral.ObjectID
	Out   string
	Zones astral.Zone
}

func (c Conn) Describe(args DescribeArgs) (r map[string]any, err error) {
	args.Out = "json"
	return rpc.Query[map[string]any](c.Conn, "describe", args)
}

type showArgs struct {
	ID astral.ObjectID
}

func (c Conn) Show(id astral.ObjectID) (r string, err error) {
	return rpc.Query[string](c.Conn, "show_object", showArgs{ID: id})
}

type ObjectResponse[T any] struct {
	Payload T        `json:"payload"`
	Type    string   `json:"type"`
	Data    []string `json:"data"`
}
