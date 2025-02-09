package apphost

import (
	"bufio"
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/query"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"github.com/cryptopunkscc/portal/runtime/rpc2/router"
	"regexp"
	"strings"
)

func Serve() cmd.Handler {
	return cmd.Handler{
		Name: "-s", Desc: "Serves rpc handler API via apphost interface.",
		Func: ServeFunc,
	}
}

func ServeFunc(ctx context.Context, root *cmd.Root) error {
	handler := cmd.Handler(*root)
	r := Default().Router(handler)
	return r.Run(ctx)
}

func (r RpcBase) Router(handler cmd.Handler, routes ...string) *Router {
	name := handler.Name
	if handler.Name != "" {
		name = strings.ReplaceAll(handler.Names()[0], "-", ".")
	}
	return &Router{
		routes: routes,
		Port:   apphost.NewPort(name),
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
	client apphost.Client
	Port   apphost.Port
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
	if r.client == nil {
		r.client = apphost.DefaultClient
	}
	r.Dependencies = append([]any{ctx}, r.Dependencies)
	routes := r.routes
	if len(routes) == 0 {
		handler := *r.Registry.Get()
		handler.Name = ""
		routes = getRoutes(nil, handler)
	}
	return r.register(ctx)
}

var RouteAll = "RouteAll"

func getRoutes(port apphost.Port, handler cmd.Handler) (r []string) {
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

func (r *Router) register(ctx context.Context) (err error) {
	listener, err := r.client.Register()
	if err != nil {
		return
	}
	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()
	var q apphost.PendingQuery
	for {
		if q, err = listener.Next(); err != nil {
			return
		}
		rr := *r
		command, ok := rr.Port.ParseCmd(q.Query())
		if !ok {
			_ = q.Close()
			continue
		}
		rr.setup(command)
		go rr.routeQuery(q)
	}
}

func (r *Router) routeQuery(q apphost.PendingQuery) {
	rr := *r
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
