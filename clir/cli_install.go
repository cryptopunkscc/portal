package clir

type Install func(string) error

func (c Cli) Install(handle Install) {
	flags := struct {
		Path string `pos:"1" default:""`
	}{}
	cmd := c.clir.NewSubCommand("i", "Install app from a given portal app bundle path.")
	cmd.AddFlags(&flags)
	cmd.Action(func() error {
		return handle(flags.Path)
	})
	return
}
