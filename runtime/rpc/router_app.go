package rpc

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
)

type App struct {
	Router
	client apphost.Client
}

func NewApp(port string) (s *App) {
	s = &App{Router: *NewRouter(port)}
	s.client = Apphost
	s.Router.RegisterRoute = s.registerRoute
	return
}

func (s *App) Client(client apphost.Client) *App {
	s.client = client
	return s
}

func (s *App) registerRoute(route string) (await func(ctx context.Context), err error) {
	listener, err := s.client.Register(route)
	if err != nil {
		return
	}
	await = func(ctx context.Context) {
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
	return
}

func (s *App) routeQuery(ctx context.Context, query apphost.QueryData) (err error) {
	if s.logger != nil {
		s.logger.Println("<~", query.Query())
	}

	// setup
	r := s.Query(query.Query())
	if r.Registry.IsEmpty() && !r.Registry.HasNext() {
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
