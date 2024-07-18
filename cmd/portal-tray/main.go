package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/runner/tray"
	. "github.com/cryptopunkscc/portal/target"
)

func main() {
	mod := Module{}
	ctx := context.Background()
	log := plog.New().Type(mod).Set(&ctx)
	mod.Conn = rpc.NewRequest(id.Anyone, PortPortal.String())
	mod.Conn.Logger(log)
	run := mod.Tray()
	err := run(ctx)
	if err != nil {
		log.P().Println(err)
	}
}

type Module struct{ portalApi }
type portalApi struct{ rpc.Conn }

func (d *Module) Tray() Tray              { return tray.NewRun(d) }
func (p portalApi) Await()                { _ = rpc.Command(p, "") }
func (p portalApi) Ping() error           { return rpc.Command(p, "ping") }
func (p portalApi) Open(src string) error { return rpc.Command(p, "open", src) }
func (p portalApi) Close() error          { return rpc.Command(p, "close") }
