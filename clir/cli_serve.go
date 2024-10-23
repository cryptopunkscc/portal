package clir

import (
	"github.com/cryptopunkscc/portal/api/target"
)

func (c Cli) Serve(handle target.Request) {
	cmd := c.clir.NewSubCommand("s", "Start portal daemon.")
	cmd.Action(func() error {
		return handle(c.ctx, "")
	})
	return
}
