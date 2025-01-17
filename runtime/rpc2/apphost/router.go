package apphost

import (
	"bufio"
	"context"
	"errors"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/query"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"github.com/cryptopunkscc/portal/runtime/rpc2/router"
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
	r := NewRouter(handler)
	return r.Run(ctx)
}

func NewRouter(handler cmd.Handler, routes ...string) *Router {
	return Rpc(Client).Router(handler, routes...)
}

func (r RpcBase) Router(handler cmd.Handler, routes ...string) *Router {
	name := handler.Name
	if handler.Name != "" {
		name = strings.ReplaceAll(handler.Names()[0], "-", ".")
	}
	return &Router{
		routes: routes,
		Port:   api.NewPort(name),
		Base: router.Base{
			Registry: router.CreateRegistry(handler),
			Unmarshalers: []caller.Unmarshaler{
				//cli.Unmarshaler{},
				//json.Unmarshaler{},
				query.Unmarshaler{},
			},
		},
		client: r.client,
	}
}

type Router struct {
	router.Base
	Logger plog.Logger
	client api.Client
	Port   api.Port
	routes []string
}

func (r *Router) Start(ctx context.Context) (err error) {
	go func() {
		if err = r.Run(ctx); err != nil {
			plog.Get(ctx).Type(r).E().Println(err)
		}
	}()
	return nil
}

func (r *Router) Run(ctx context.Context) error {
	r.Dependencies = append([]any{ctx}, r.Dependencies)
	routes := r.routes
	if len(routes) == 0 {
		handler := *r.Registry.Get()
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

var RouteAll = "RouteAll"

func getRoutes(port api.Port, handler cmd.Handler) (r []string) {
	if name := handler.Names()[0]; name != "" {
		port = port.Add(name)
	}
	if handler.Func != nil {
		p := port.String()
		if handler.Func == RouteAll {
			p += "*"
		}
		r = append(r, p)
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
	listener, err := r.client.Register(port)
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
	defer conn.Close()
	flow := NewClient(conn)
	if r.Logger != nil {
		flow.Logger(r.Logger.Scope(q.Query()))
	}
	r.Dependencies = append(r.Dependencies, flow)
	rr.Dependencies = r.Dependencies
	scanner := bufio.NewScanner(conn)
	for {
		if !rr.skip() {
			if err = rr.Respond(flow.Serializer); err != nil {
				return
			}
		} else {
			rrr := rr
			r = &rrr
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

var ErrUnauthorized = errors.New("unauthorized")

func (r *Router) skip() bool {
	return r.Registry.Get().Func == RouteAll
}

func (r *Router) setup(query string) {
	r.Base = r.Query(query)
}

func (r *Router) authorize() bool {
	rr := r.Query("!")
	if rr.Registry != r.Registry {
		return <-rr.Call() != false
	}
	return true
}
