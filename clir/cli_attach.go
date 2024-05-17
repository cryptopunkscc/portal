package clir

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func (c Cli) Attach(handler target.Spawn) {
	flags := &struct {
		Runner   string `pos:"1" description:"App runner [goja, wails]."`
		Absolute string `pos:"2" description:"Absolute path to app bundle or directory."`
	}{}
	cmd := c.clir.NewSubCommand("r", "Start portal app in given runner.")
	cmd.AddFlags(flags)
	cmd.Action(func() (err error) {
		return handler(c.ctx, flags.Absolute)
	})
	return
}
