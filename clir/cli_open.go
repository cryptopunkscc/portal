package clir

import (
	"github.com/cryptopunkscc/portal/target"
)

func (c Cli) Open(handler target.Dispatch) {
	flags := &struct {
		Absolute string `pos:"1" description:"Absolute path to app bundle or directory."`
	}{}
	cmd := c.clir.NewSubCommand("o", "Start portal app in given runner.")
	cmd.AddFlags(flags)
	cmd.Action(func() (err error) {
		return handler(c.ctx, flags.Absolute)
	})
	return
}
