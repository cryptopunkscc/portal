package msg

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"strings"
	"time"
)

type Handler struct {
	reloader Reloader
	apphost  target.ApphostCache
	changes  mem.Cache[time.Time]
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
		changes:  mem.NewCache[time.Time](),
	}
}

func (s *Handler) HandleMsg(ctx context.Context, msg target.Msg) {
	log := plog.Get(ctx).D()
	log.Println("received broadcast message:", msg)
	switch msg.Event {
	case target.DevChanged:
		for _, c := range s.apphost.Connections() {
			if c.In {
				continue
			}
			query := strings.TrimPrefix(c.Query, "dev.")
			if strings.HasPrefix(query, msg.Pkg) {
				s.changes.Set(msg.Pkg, msg.Time)
			}
		}
	case target.DevRefreshed:
		if ok := s.changes.Delete(msg.Pkg); !ok || s.changes.Size() > 0 {
			log.Println("cannot reload", ok, s.changes.Size())
			return
		}
		log.Println("reloading")
		if err := s.reloader.Reload(); err != nil {
			log.F().Println(err)
		}
	}
}
