package clir

import (
	"github.com/cryptopunkscc/portal/target"
)

func (c Cli) Start(handle target.Request) {
	flags := &struct {
		Src  string `pos:"1" default:"" description:"Application source. The source can be a app name, package name, app bundle path or app dir."`
		Type string `pos:"2" default:"2"`
	}{}
	f := func() error {
		return handle(c.ctx, flags.Src, flags.Type)
	}
	cmd := c.clir.NewSubCommand(
		"d",
		"Start a given source to be run as an application(s). The source can be an app name, package name, app bundle path, or a directory containing application(s).",
	)
	cmd.AddFlags(flags)
	cmd.Action(f)
	c.clir.DefaultCommand(cmd)
	c.clir.AddFlags(flags)
	c.clir.Action(f)
	return
}
