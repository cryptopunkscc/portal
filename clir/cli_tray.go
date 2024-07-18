package clir

import "github.com/cryptopunkscc/portal/target"

func (c Cli) Tray(handle target.Tray) {
	f := func() error { return handle(c.ctx) }
	cmd := c.clir.NewSubCommand("t", "Run tray indicator for portal-app")
	cmd.Action(f)
	return
}
