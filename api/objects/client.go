package objects

import (
	"bytes"
	"encoding/binary"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/objects"
	"github.com/cryptopunkscc/astrald/object"
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
	ID     object.ID
	Offset astral.Uint64
	Zone   astral.Zone
}

func (c Conn) Read(args ReadArgs) (r io.ReadCloser, err error) {
	return rpc.NewCall(c.Conn, "read", args)
}

func (c Conn) Search(args SearchArgs) (out <-chan ObjectResponse[objects.SearchResult], err error) {
	args.Format = "json"
	return rpc.Subscribe[ObjectResponse[objects.SearchResult]](c.Conn, "search", args)
}

type ObjectResponse[T any] struct {
	Payload T        `json:"payload"`
	Type    string   `json:"type"`
	Data    []string `json:"data"`
}

type SearchArgs struct {
	Query  string `query:"q"`
	Zone   astral.Zone
	Format string
	Ext    string
}

func (c Conn) Describe(args DescribeArgs) (r map[string]any, err error) {
	args.Format = "json"
	return rpc.Query[map[string]any](c.Conn, "describe", args)
}

type DescribeArgs struct {
	ID     object.ID
	Format string
	Zones  astral.Zone
}

func (c Conn) Show(id object.ID) (r string, err error) {
	return rpc.Query[string](c.Conn, "show_object", showArgs{ID: id})
}

type showArgs struct {
	ID object.ID
}
