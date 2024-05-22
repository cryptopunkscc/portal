package clir

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

func (c Cli) Dev(handle target.Dispatch) {
	flags := &struct {
		Src string `pos:"1" default:"" description:"Application source. The source can be a app name, package name, app bundle or app dir."`
	}{}
	f := func() error {
		return handle(c.ctx, flags.Src)
	}
	cmd := c.clir.NewSubCommand("d", "Dispatch project or app from a given source in development environment.")
	cmd.AddFlags(flags)
	cmd.Action(f)
	c.clir.DefaultCommand(cmd)
	c.clir.AddFlags(flags)
	c.clir.Action(f)
	return
}
