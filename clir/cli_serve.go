package clir

import (
	"context"
)

type Serve func(
	ctx context.Context,
) error

func (c Cli) Serve(handle Serve) {
	cmd := c.clir.NewSubCommand("s", "Start portal daemon.")
	cmd.Action(func() error {
		return handle(c.ctx)
	})
	return
}
