package rpc

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/net"
	"github.com/cryptopunkscc/astrald/node"
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

type Module struct {
	rpc.Router
	node node.Node
}

func NewModule(node node.Node, port string) (r *Module) {
	r = &Module{Router: *rpc.NewRouter(port), node: node}
	r.Router.RegisterRoute = r.registerRoute
	return
}

func (m Module) registerRoute(route string) (await func(ctx context.Context), err error) {
	if err = m.node.LocalRouter().AddRoute(route, m); err != nil {
		return
	}
	await = func(ctx context.Context) {
		<-ctx.Done()
		_ = m.node.LocalRouter().RemoveRoute(route)
	}
	return
}

func (m Module) RouteQuery(ctx context.Context, query net.Query, caller net.SecureWriteCloser, hints net.Hints) (s net.SecureWriteCloser, err error) {
	// setup
	m.Router = *m.Query(query.Query())
	if m.Registry.IsEmpty() && query.Query() != m.Port {
		return nil, net.ErrRejected
	}

	// authorize
	if m.authorize(ctx, query.Caller(), query) {
		return nil, net.ErrRejected
	}

	// accept
	return net.Accept(query, caller, func(conn net.SecureConn) {
		_ = m.Handle(ctx, query, query.Caller(), conn)
	})
}

func (m Module) authorize(ctx context.Context, callerID id.Identity, query any) bool {
	res, _ := m.Command("!").With(ctx, query).Call()
	if len(res) > 0 {
		switch v := res[0].(type) {
		case bool:
			return v
		case string:
			return m.node.Auth().Authorize(callerID, v)
		}
	}
	return false
}
