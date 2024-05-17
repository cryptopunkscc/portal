package clir

type Build func(string) error

func (c Cli) Build(handle Build) {
	flags := struct {
		Path string `pos:"1" default:"."`
	}{}
	cmd := c.clir.NewSubCommand("b", "Build project and generate portal app bundle.")
	cmd.AddFlags(&flags)
	cmd.Action(func() (err error) {
		return handle(flags.Path)
	})
	return
}
