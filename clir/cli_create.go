package clir

import "strings"

type Templates func() error

type Create func(
	projectName string,
	targetDir string,
	templates []string,
	force bool,
) (err error)

func (c Cli) Create(
	templates Templates,
	create Create,
) {
	emptyFlags := struct {
		Dir      string `pos:"1" description:"Project directory"`
		Name     string `name:"n" description:"Name of project"`
		Template string `name:"t" description:"Name of built-in template to use, path to template or template url"`
		Force    bool   `name:"f" description:"Force recreate project"`
		List     bool   `name:"l" description:"List available templates"`
	}{}
	flags := emptyFlags
	cmd := c.clir.NewSubCommand("c", "Create new project from template.")
	cmd.AddFlags(&flags)
	cmd.Action(func() error {
		switch {
		case flags == emptyFlags || flags.List:
			return templates()
		default:
			temps := strings.Split(flags.Template, " ")
			return create(flags.Name, flags.Dir, temps, flags.Force)
		}
	})
	return
}
