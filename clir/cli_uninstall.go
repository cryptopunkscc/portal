package clir

import "context"

type Uninstall func(context.Context, string) error

func (c Cli) Uninstall(handle Uninstall) {
	flags := struct {
		Id string `pos:"1" default:""`
	}{}
	cmd := c.clir.NewSubCommand("u", "Uninstall app by given id. The id can by app name or package name")
	cmd.AddFlags(&flags)
	cmd.Action(func() error {
		return handle(c.ctx, flags.Id)
	})
	return
}
