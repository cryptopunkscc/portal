package clir

import (
	"github.com/cryptopunkscc/portal/api/target"
)

func (c Cli) Start(handle target.Request) {
	flags := &struct {
		Src string `pos:"1" default:"" description:"Application source. The source can be a app name, package name, app bundle path or app dir."`
	}{}
	f := func() error {
		args := c.clir.OtherArgs()
		if len(args) > 0 {
			args = args[1:]
		}
		return handle(c.ctx, flags.Src, args...)
	}
	c.clir.AddFlags(flags)
	c.clir.Action(f)
	return
}
