package clir

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

func (c Cli) Dispatch(handle target.Dispatch) {
	flags := &struct {
		Src string `pos:"1" default:""`
	}{}
	f := func() error {
		return handle(c.ctx, flags.Src)
	}
	cmd := c.clir.NewSubCommand(
		"d",
		"Dispatch given source to be run as application. The source can be a app name, package name, app bundle path or app dir.",
	)
	cmd.AddFlags(flags)
	cmd.Action(f)
	c.clir.DefaultCommand(cmd)
	c.clir.AddFlags(flags)
	c.clir.Action(f)
	return
}
