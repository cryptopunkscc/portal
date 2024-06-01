package msg

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/registry"
	"github.com/cryptopunkscc/go-astral-js/target"
	"strings"
	"time"
)

type Handler struct {
	reloader Reloader
	apphost  target.ApphostCache
	changes  registry.Cache[time.Time]
}

type Reloader interface {
	Reload() error
}

func NewHandler(
	reloader Reloader,
	apphost target.ApphostCache,
) *Handler {
	return &Handler{
		reloader: reloader,
		apphost:  apphost,
		changes:  registry.New[time.Time](),
	}
}

func (s *Handler) HandleMsg(ctx context.Context, msg target.Msg) {
	log := plog.Get(ctx)
	log.Println("received ctrl message:", msg)
	switch msg.Event {
	case target.DevChanged:
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
	case target.DevRefreshed:
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
