package clir

import "context"

type Build func(context.Context, string) error

func (c Cli) Build(handle Build) {
	flags := struct {
		Path string `pos:"1" default:"."`
	}{}
	cmd := c.clir.NewSubCommand("b", "Build project and generate portal app bundle.")
	cmd.AddFlags(&flags)
	cmd.Action(func() (err error) {
		return handle(c.ctx, flags.Path)
	})
	return
}
