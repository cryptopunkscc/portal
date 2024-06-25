package clir

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

func (c Cli) Dev(handle target.Dispatch) {
	flags := &struct {
		Src  string `pos:"1" default:"." description:"Application source. The source can be a app name, package name, app bundle path or app dir."`
		Type string `pos:"2" default:"2"`
	}{}
	f := func() error {
		return handle(c.ctx, flags.Src, flags.Type)
	}
	cmd := c.clir.NewSubCommand("d", "Dispatch project or app from a given source in development environment.")
	cmd.AddFlags(flags)
	cmd.Action(f)
	c.clir.DefaultCommand(cmd)
	c.clir.AddFlags(flags)
	c.clir.Action(f)
	return
}
