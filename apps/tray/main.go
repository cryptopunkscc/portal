package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/portal"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/client"
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

func main() {
	mod := module{}
	ctx := context.Background()
	log := plog.New().Type(mod).Set(&ctx)
	conn := rpc.NewRequest(id.Anyone, PortPortal.String())
	conn.Logger(log)
	mod.Client = client.PortalClient{Conn: conn}
	r := newRunner(mod)
	err := r.Run(ctx)
	if err != nil {
		log.P().Println(err)
	}
}

type module struct{ portal.Client }
