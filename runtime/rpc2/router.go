package rpc

import (
	"fmt"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"github.com/cryptopunkscc/portal/runtime/rpc2/registry"
	"strings"
)

type Router struct {
	Registry     *registry.Node[*cmd.Handler]
	Unmarshalers []caller.Unmarshaler
	Dependencies []any
	query        string
	args         string
	err          error
}

func CreateRegistry(handler cmd.Handler) *registry.Node[*cmd.Handler] {
	handler.Name = ""
	r := registry.New[*cmd.Handler]('.', ' ')
	InjectHandler(r, handler)
	return r
}

func InjectHandler(registry *registry.Node[*cmd.Handler], handler cmd.Handler) {
	for _, name := range handler.Names() {
		r := registry.Add(name, &handler)
		for _, h := range handler.Sub {
			InjectHandler(r, h)
		}
	}
}

func (r Router) Query(query string) Router {
	r.Registry, r.args = r.Registry.Fold(query)
	if r.Registry.IsEmpty() {
		r.err = fmt.Errorf("invalid query: %s", query)
		return r
	}
	return r
}

func (r Router) Call() (result any) {
	if r.err != nil {
		return r.err
	}
	in := []byte(r.args)
	out, err := r.caller().
		Unmarshalers(r.Unmarshalers...).
		Defaults(r).
		Defaults(r.Dependencies...).
		Call(in)
	switch {
	case err != nil:
		return err
	case len(out) == 0:
		return nil
	case len(out) == 1:
		return out[0]
	default:
		return out
	}
}

func (r Router) caller() *caller.Func {
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
