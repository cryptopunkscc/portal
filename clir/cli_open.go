package clir

import (
	"github.com/cryptopunkscc/portal/target"
)

func (c Cli) Open(handle target.Request) {
	flags := &struct {
		Absolute string `pos:"1" description:"Absolute path to app bundle or directory."`
	}{}
	cmd := c.clir.NewSubCommand("o", "Start portal app in given runner.")
	cmd.AddFlags(flags)
	cmd.Action(func() (err error) {
		return handle(c.ctx, flags.Absolute)
	})
	return
}
