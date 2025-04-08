package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"strings"
	"time"
)

type handler struct {
	reload  Reload
	cache   apphost.Cache
	changes mem.Cache[time.Time]
}

type Reload func() error

func newHandler(
	reload Reload,
	cache apphost.Cache,
) *handler {
	return &handler{
		reload:  reload,
		cache:   cache,
		changes: mem.NewCache[time.Time](),
	}
}

func (s *handler) HandleMsg(ctx context.Context, msg target.Msg) {
	log := plog.Get(ctx).D()
	log.Println("received broadcast message:", msg)
	switch msg.Event {
	case target.DevChanged:
		if s.cache == nil {
			return
		}
		for _, c := range s.cache.Connections().Copy() {
			if c.In() {
				continue
			}
			query := strings.TrimPrefix(c.Query(), "dev.")
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
		if err := s.reload(); err != nil {
			log.F().Println(err)
		}
	}
}
