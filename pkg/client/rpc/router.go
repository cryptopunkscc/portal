package rpc

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
	"github.com/cryptopunkscc/portal/pkg/util/rpc"
	"github.com/cryptopunkscc/portal/pkg/util/rpc/caller/query"
	"github.com/cryptopunkscc/portal/pkg/util/rpc/cmd"
	"github.com/cryptopunkscc/portal/pkg/util/rpc/router"
)

func (r Rpc) Router(handler cmd.Handler) rpc.Router {
	rr := &Router{}
	rr.Rpc = r
	rr.Registry = router.CreateRegistry(handler)
	rr.Unmarshal = query.Unmarshal
	rr.Register = r.Register
	return rr
}

type Router struct {
	Rpc
	router.Base
	Register func(ctx context.Context) (*astrald.Listener, error)
	ctx      context.Context
	listener *astrald.Listener
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

	r.listener, err = r.Register(ctx)
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
	var q *astrald.PendingQuery
	for {
		if q, err = r.listener.Next(); err != nil {
			if errors.Is(r.ctx.Err(), context.Canceled) {
				err = nil
			}
			return
		}
		rr := *r
		qq := *q
		go func() {
			if err := rr.routeQuery(&qq); err != nil && r.Log != nil {
				r.Log.E().Println(qq)
			}
		}()
	}
}

type PendingQuery interface {
	Query() string
	Reject() error
	Accept() *apphost.Conn
}

func (r *Router) routeQuery(q PendingQuery) (err error) {
	r.Responses = math.MaxInt64
	r.Add(&r.Base, q)

	rr := *r
	command := q.Query()
	rr.Setup(command)
	if !rr.authorize() {
		_ = q.Reject()
		return
	}
	conn := q.Accept()
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
