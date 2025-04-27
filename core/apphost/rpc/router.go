package rpc

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/caller/query"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/pkg/rpc/router"
	"math"
	"strings"
)

func (r *Rpc) Router(handler cmd.Handler) *Router {
	rr := &Router{}
	rr.Rpc = *r
	rr.Registry = router.CreateRegistry(handler)
	rr.Unmarshal = query.Unmarshal
	return rr
}

type Router struct {
	Rpc
	router.Base
	ctx      context.Context
	listener apphost.Listener
}

func (r *Router) Start(ctx context.Context) (err error) {
	if err = r.Init(ctx); err != nil {
		return
	}
	go func() {
		if err := r.Listen(); err != nil {
			plog.Get(ctx).Type(r).E().Println(err)
		}
	}()
	return
}

func (r *Router) Run(ctx context.Context) (err error) {
	if err = r.Init(ctx); err != nil {
		return
	}
	return r.Listen()
}

func (r *Router) Init(ctx context.Context) (err error) {
	defer plog.TraceErr(&err)
	r.ctx = ctx
	r.Dependencies = append([]any{ctx}, r.Dependencies...)

	r.listener, err = r.Apphost.Register()
	if err != nil {
		return
	}
	go func() {
		<-ctx.Done()
		_ = r.listener.Close()
	}()
	return
}

func (r *Router) Listen() (err error) {
	var q apphost.PendingQuery
	for {
		if q, err = r.listener.Next(); err != nil {
			if errors.Is(r.ctx.Err(), context.Canceled) {
				err = nil
			}
			return
		}
		rr := *r
		go func() {
			if err := rr.routeQuery(q); err != nil && r.Log != nil {
				r.Log.E().Println(q)
			}
		}()
	}
}

func (r *Router) routeQuery(q apphost.PendingQuery) (err error) {
	r.Responses = math.MaxInt64
	r.Add(&r.Base, q)

	rr := *r
	command := q.Query()
	rr.Setup(command)
	if !rr.authorize() {
		_ = q.Reject()
		return
	}
	conn, err := q.Accept()
	if err != nil {
		return
	}
	defer conn.Close()

	s := r.client(conn)
	if r.Log != nil {
		s.Logger(r.Log.Scope(q.Query()))
	}

	r.Add(s)
	rr.Add(s)

	for {
		if !rr.IsEmpty() {
			if err = rr.Respond(s.Serializer); err != nil {
				return
			}
		}
		if r.Responses == 0 {
			return
		}
		r.Responses--
		command, err = s.ReadString('\n')
		if err != nil {
			return
		}
		command = strings.TrimSpace(command)
		rr = *r
		rr.Setup(command)
		if !rr.authorize() {
			_ = s.Encode(ErrUnauthorized)
			return
		}
	}
}

var ErrUnauthorized = errors.New("unauthorized")

func (r *Router) authorize() bool {
	rr := r.Query("!")
	if rr.Registry != r.Registry {
		return <-rr.Call() != false
	}
	return true
}
