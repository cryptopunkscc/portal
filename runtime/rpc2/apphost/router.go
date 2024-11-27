package apphost

import (
	"bufio"
	"context"
	"errors"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/port"
	runtime "github.com/cryptopunkscc/portal/runtime/apphost"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/json"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/query"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"regexp"
	"strings"
	"sync"
)

func Serve() cmd.Handler {
	return cmd.Handler{
		Name: "-s", Desc: "Serves rpc handler API via apphost interface.",
		Func: ServeFunc,
	}
}

func ServeFunc(ctx context.Context, root *cmd.Root) error {
	handler := cmd.Handler(*root)
	r := NewRouter(handler, nil)
	return r.Run(ctx)
}

var Client api.Client = runtime.Default()

type Router struct {
	rpc.Router
	Port   port.Port
	routes []string
}

var ErrUnauthorized = errors.New("unauthorized")

func NewRouter(handler cmd.Handler, port port.Port, routes ...string) *Router {
	if len(port) == 0 && handler.Name != "" {
		name := strings.ReplaceAll(handler.Names()[0], "-", ".")
		port = port.Add(name)
	}
	return &Router{
		routes: routes,
		Port:   port,
		Router: rpc.Router{
			Registry: rpc.CreateRegistry(handler),
			Unmarshalers: []caller.Unmarshaler{
				cli.Unmarshaler{},
				json.Unmarshaler{},
				query.Unmarshaler{},
			},
		},
	}
}

func (r *Router) Run(ctx context.Context) error {
	routes := r.routes
	if len(routes) == 0 {
		handler := *r.Router.Registry.Get()
		handler.Name = ""
		routes = getRoutes(nil, handler)
	}

	wg := sync.WaitGroup{}
	errs := make(chan error, len(routes))
	wg.Add(len(routes))
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

func getRoutes(port port.Port, handler cmd.Handler) (r []string) {
	if name := handler.Names()[0]; name != "" {
		port = port.Add(name)
	}
	if handler.Func != nil {
		r = append(r, port.String())
	}
	for _, h := range handler.Sub {
		if b, _ := regexp.MatchString(`^[a-z]+`, h.Name); !b {
			continue
		}
		if strings.HasPrefix(h.Name, "help") {
			continue
		}
		r = append(r, getRoutes(port, h)...)
	}
	return
}

func (r *Router) formatPort(route string) (port string) {
	switch route {
	case "*":
		port = r.Port.String() + "*"
	default:
		port = r.Port.Add(route).String()
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
	command := rr.Port.ParseCmd(q.Query())
	rr.setup(command)
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
		if err = rr.Respond(flow.Serializer); err != nil {
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
	return <-r.Router.Query("!").Call() != false
}
