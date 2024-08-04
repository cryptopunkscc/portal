package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	. "github.com/cryptopunkscc/portal/target"
)

func main() {
	mod := module{}
	ctx := context.Background()
	log := plog.New().Type(mod).Set(&ctx)
	mod.Conn = rpc.NewRequest(id.Anyone, PortPortal.String())
	mod.Conn.Logger(log)
	r := newRunner(mod.portalApi)
	err := r.Run(ctx)
	if err != nil {
		log.P().Println(err)
	}
}

type module struct{ portalApi }
