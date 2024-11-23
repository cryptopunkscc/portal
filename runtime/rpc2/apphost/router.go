package apphost

import (
	"bufio"
	"context"
	"errors"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	runtime "github.com/cryptopunkscc/portal/runtime/apphost"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"golang.org/x/exp/maps"
	"sync"
)

var Client api.Client = runtime.Default()

type Router struct {
	rpc.Router
	target.Port
	routes []string
}

var ErrUnauthorized = errors.New("unauthorized")

func NewRouter(router rpc.Router, port target.Port, routes ...string) *Router {
	return &Router{Router: router, Port: port, routes: routes}
}

func (r *Router) Run(ctx context.Context) error {
	routes := r.routes
	if len(routes) == 0 {
		routes = maps.Keys(r.Router.Registry.All())
	}

	wg := sync.WaitGroup{}
	errs := make(chan error, len(routes))
	for _, route := range routes {
		port := r.formatPort(route)
		go func() {
			defer wg.Done()
			if err := r.register(ctx, port); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)
	var errsArr []error
	for e := range errs {
		errsArr = append(errsArr, e)
	}
	return errors.Join(errsArr...)
}

func (r *Router) formatPort(route string) (port string) {
	switch route {
	case "*":
		port = r.Port.String() + "*"
	default:
		port = r.Port.Route(route).String()
	}
	return
}

func (r *Router) register(ctx context.Context, port string) (err error) {
	listener, err := Client.Register(port)
	if err != nil {
		return
	}
	queries := listener.QueryCh()
	for {
		select {
		case q := <-queries:
			rr := *r
			go rr.routeQuery(q)
		case <-ctx.Done():
			if err = listener.Close(); err != nil {
				return
			}
			return
		}
	}
}

func (r *Router) routeQuery(q api.QueryData) {
	rr := *r
	rr.setup(q.Query())
	if !rr.authorize() {
		_ = q.Reject()
		return
	}
	conn, err := q.Accept()
	if err != nil {
		return
	}
	flow := NewClient(conn)
	scanner := bufio.NewScanner(conn)
	for {
		result := rr.Router.Call()
		if err = flow.Encode(result); err != nil {
			return
		}
		if !scanner.Scan() {
			return
		}
		text := scanner.Text()
		rr = *r
		rr.setup(text)
		if !rr.authorize() {
			_ = flow.Encode(ErrUnauthorized)
			return
		}
	}
}

func (r *Router) setup(query string) {
	r.Router = r.Router.Query(query)
}

func (r *Router) authorize() bool {
	return r.Router.Query("!").Call() != false
}
