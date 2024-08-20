package mobile

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/apphost"
	. "github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
	"strings"
)

type Deps interface {
	Api() Api
	Ctx() context.Context
	Port() target.Port
	Serve() target.Request
	Cache() apphost.Cache
}

func Feat(deps Deps) Serve {
	api := deps.Api()
	ctx := deps.Ctx()
	port := deps.Port()
	serve := deps.Serve()
	events := deps.Cache().Events()

	return func() {
		go subscribe(ctx, api, port, events)
		go start(ctx, api, serve)
	}
}

func subscribe(
	ctx context.Context,
	api Api,
	port target.Port,
	events *sig.Queue[apphost.Event],
) {
	all := 5 // TODO resolve handlers count automatically
	connected := 0
	for event := range events.Subscribe(ctx) {
		if strings.HasPrefix(event.Port, port.String()) {
			switch event.Type {
			case apphost.Register:
				connected++
			case apphost.Unregister:
				connected--
			default:
				continue
			}
			switch connected {
			case all:
				api.Event(&Event{Msg: STARTED})
			}
		}
	}
}

func start(ctx context.Context, api Api, serve target.Request) {
	api.Event(&Event{Msg: STARTING})
	err := serve(ctx, "")
	api.Event(&Event{Msg: STOPPED, Err: err})
}
