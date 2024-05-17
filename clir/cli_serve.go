package clir

import (
	"context"
)

type Serve func(
	ctx context.Context,
	tray bool,
) error

func (c Cli) Serve(handle Serve) {
	flags := &struct {
		Tray bool `name:"t" description:"Launch tray indicator."`
	}{}
	cmd := c.clir.NewSubCommand("s", "Start portal daemon.")
	cmd.AddFlags(flags)
	cmd.Action(func() error {
		return handle(c.ctx, flags.Tray)
	})
	return
}
