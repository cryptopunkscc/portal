package router

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/core/rpc"
	"github.com/cryptopunkscc/portal/core/rpc/caller"
	"github.com/cryptopunkscc/portal/core/rpc/cmd"
	"github.com/cryptopunkscc/portal/core/rpc/registry"
	"github.com/cryptopunkscc/portal/core/rpc/stream"
	"io"
	"strings"
)

type Base struct {
	Registry     *registry.Node[*cmd.Handler]
	Unmarshal    caller.Unmarshal
	Responses    int64
	Dependencies []any
	args         string
	err          error
}

func CreateRegistry(handler cmd.Handler) *registry.Node[*cmd.Handler] {
	r := registry.New[*cmd.Handler]('.', ' ')
	rr := r.Add("", &handler)
	for _, h := range handler.Sub {
		injectHandler(rr, h)
	}
	return r
}

func injectHandler(registry *registry.Node[*cmd.Handler], handler cmd.Handler) {
	for _, name := range handler.Names() {
		r := registry.Add(name, &handler)
		for _, h := range handler.Sub {
			injectHandler(r, h)
		}
	}
}

func (r *Base) IsEmpty() bool {
	return !(!r.Registry.IsEmpty() && r.Registry.Get().Func != nil)
}

func (r *Base) Add(dependencies ...any) *Base {
	r.Dependencies = append(r.Dependencies, dependencies...)
	return r
}

func (r *Base) Err() error { return r.err }

func (r *Base) Setup(query string) {
	base := r.Query(query)
	if !base.IsEmpty() {
		*r = base
	}
}

func (r *Base) Query(query string) Base {
	rr := *r
	rr.Registry, rr.args = rr.Registry.Fold(query)
	if rr.Registry.IsEmpty() {
		rr.err = fmt.Errorf("invalid query: %s", query)
		return rr
	}
	rr.args = strings.TrimPrefix(rr.args, "?")
	return rr
}

func (r *Base) Respond(conn *stream.Serializer) (err error) {
	for item := range r.Call() {
		if item == rpc.Close {
			_ = conn.Close()
			return rpc.Close
		}
		if err = conn.Encode(item); err != nil {
			return
		}
	}
	//err = conn.Encode(End) //Fixme interferes with test/rpc/go_test.go and cmd/cli
	return
}

func (r *Base) Call() (o <-chan any) {
	c := make(chan any, 1)
	o = c
	if r.err != nil {
		c <- r.err
		close(c)
		return c
	}
	in := []byte(r.args)
	out, err := r.caller().
		Unmarshaler(r.Unmarshal).
		Defaults(r.Dependencies...).
		Defaults(r.Registry.Get()).
		Call(in)

	go respond(c, err, out...)
	return o
}

func respond(c chan any, err error, out ...any) {
	defer close(c)
	switch {
	case err != nil:
		if errors.Is(err, io.EOF) {
			return
		}
		c <- err
	case len(out) == 0:
		c <- stream.End
	case len(out) == 1:
		r := out[0]
		switch v := r.(type) {
		case <-chan any:
			for a := range v {
				c <- a
			}
		default:
			c <- r
		}
	default:
		c <- out
	}
}

func (r *Base) caller() *caller.Func {
	h := r.Registry.Get()
	switch v := h.Func.(type) {
	case *caller.Func:
		return v
	default:
		f := caller.New(h.Func)
		f.Names = strings.Split(h.Name, " ")
		h.Func = f
		return f
	}
}
