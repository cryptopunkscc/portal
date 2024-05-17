package clir

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
)

func (c Cli) Open(handle runtime.Spawn) {
	flags := &struct {
		Src string `pos:"1" default:""`
	}{}
	f := func() error {
		return handle(c.ctx, flags.Src)
	}
	cmd := c.clir.NewSubCommand("o", "Open app from a given source. The source can be a app name, package name, app bundle or app dir.")
	cmd.AddFlags(flags)
	cmd.Action(f)
	c.clir.DefaultCommand(cmd)
	c.clir.AddFlags(flags)
	c.clir.Action(f)
	return
}
