package objects

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/objects"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Op(rpc rpc.Rpc, target ...string) OpClient {
	if len(target) == 0 {
		target = append(target, "localnode")
	}
	return OpClient{rpc.Format("json").Request(target[0], "objects")} // change node ID to call foreign node
}

type OpClient struct{ rpc.Conn }

func (c OpClient) Push(obj astral.Object) (ok bool, err error) {
	defer plog.TraceErr(&err)
	buf := bytes.NewBuffer(nil)
	if _, err = astral.WriteCanonical(buf, obj); err != nil {
		return
	}
	conn := c.Copy()
	if err = rpc.Call(conn, "push", rpc.Opt{"size": buf.Len()}); err != nil {
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

func (c OpClient) Read(args ReadArgs) (r io.ReadCloser, err error) {
	return rpc.NewCall(c.Conn, "read", args)
}

func (c OpClient) Fetch(args ReadArgs, obj astral.Object) (err error) {
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
	Query string `query:"q" include:"empty"`
	Zone  astral.Zone
	Out   string
	Ext   string // not implemented yet
}

func (c OpClient) Search(args SearchArgs) (out <-chan rpc.Json[objects.SearchResult], err error) {
	args.Out = "json"
	return rpc.Subscribe[rpc.Json[objects.SearchResult]](c.Conn, "search", args)
}

type ScanArgs struct {
	Type   string
	Repo   string
	Out    string
	Follow bool
	Zone   astral.Zone
}

func (c OpClient) Scan(args ScanArgs) (out <-chan astral.ObjectID, err error) {
	args.Out = "json"
	o, err := rpc.Subscribe[rpc.Json[astral.ObjectID]](c.Conn, "scan", args)
	if err != nil {
		return
	}
	out = flow.Map(o, func(t1 rpc.Json[astral.ObjectID]) (astral.ObjectID, bool) { return t1.Object, true })
	return
}

func (c OpClient) Scan2(args ScanArgs) (out <-chan rpc.Json[string], err error) {
	args.Out = "json"
	return rpc.Subscribe[rpc.Json[string]](c.Conn, "scan", args)
}

type DescribeArgs struct {
	ID    astral.ObjectID
	Out   string
	Zones astral.Zone
}

func (c OpClient) Describe(args DescribeArgs) (r <-chan map[string]any, err error) {
	args.Out = "json"
	return rpc.Subscribe[map[string]any](c.Conn, "describe", args)
}

func (c OpClient) Show(id astral.ObjectID) (r string, err error) {
	return rpc.Query[string](c.Conn, "show_object", rpc.Opt{"id": id})
}
