package rpc

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/astral"
)

type App struct {
	Router
}

func NewApp(port string) (s *App) {
	s = &App{Router: *NewRouter(port)}
	s.Router.registerRoute = s.registerRoute
	return
}

func (s *App) registerRoute(ctx context.Context, route string) (err error) {
	listener, err := astral.Register(route)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	done := ctx.Done()
	queries := listener.QueryCh()
	for {
		select {
		case <-done:
			return
		case q := <-queries:
			ss := *s
			go ss.routeQuery(ctx, q)
		}
	}
}

func (s *App) routeQuery(ctx context.Context, query *astral.QueryData) (err error) {
	s.logger.Println("<", query.Query())

	// setup
	r := s.Query(query.Query())
	if r.registry.IsEmpty() && query.Query() != r.port {
		return query.Reject()
	}

	// authorize
	if !r.Authorize(ctx, query) {
		return query.Reject()
	}

	// accept
	conn, err := query.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	_ = r.Handle(ctx, query, query.RemoteIdentity(), conn)
	return
}
