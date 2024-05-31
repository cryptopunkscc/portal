package service_dev

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/broadcast"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/registry"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
	"strings"
	"time"
)

type Service struct {
	reloader Reloader
	apphost  target.ApphostCache
	changes  registry.Cache[time.Time]
}

type Reloader interface {
	Reload() error
}

func NewService(runner Reloader, apphost target.ApphostCache) *Service {
	return &Service{
		reloader: runner,
		apphost:  apphost,
		changes:  registry.New[time.Time](),
	}
}

func (s *Service) Start(ctx context.Context, portal target.Portal) {
	go func() {
		if err := s.Run(ctx, portal); err != nil {
			port := target.DevPort(portal)
			plog.Get(ctx).Type(s).Printf("%s: %v", port, err)
		}
	}()
}

func (s *Service) Run(ctx context.Context, portal target.Portal) error {
	plog.Get(ctx).Type(s).Set(&ctx)
	port := target.DevPort(portal)
	app := rpc.NewApp(port)
	app.RouteFunc("ctrl", s.handleCtrl)
	return app.Run(ctx)
}

func (s *Service) handleCtrl(ctx context.Context, msg broadcast.Msg) {
	log := plog.Get(ctx)
	log.Println("received ctrl message:", msg)
	switch msg.Event {
	case broadcast.Changed:
		log.Println(s.apphost.Connections())
		for _, c := range s.apphost.Connections() {
			if c.In {
				continue
			}
			query := strings.TrimPrefix(c.Query, "dev.")
			log.Println(query, msg.Pkg, strings.HasPrefix(query, msg.Pkg))
			if strings.HasPrefix(query, msg.Pkg) {
				s.changes.Set(msg.Pkg, msg.Time)
			}
		}
	case broadcast.Refreshed:
		log.Println(s.changes.Copy(), s.apphost.Connections(), s.changes.Size(), s.changes.Copy())

		if ok := s.changes.Delete(msg.Pkg); !ok || s.changes.Size() > 0 {
			log.Println("cannot reload", ok, s.changes.Size())
			return
		}
		log.Println("reloading")
		if err := s.reloader.Reload(); err != nil {
			plog.Get(ctx).F().Println(err)
		}
	}
}
