package clir

import (
	"github.com/cryptopunkscc/portal/feat/create"
	"strings"
)

type Templates func() error

func (c Cli) Create(
	templates Templates,
	create create.Run,
) {
	emptyFlags := struct {
		Targets string `pos:"1" description:"List of templates with optional module names like: 'svelte backend' or 'svelte:front backend:back'."`
		Dir     string `pos:"2" description:"Project directory."`
		List    bool   `name:"l" description:"List available templates."`
	}{}
	flags := emptyFlags
	c.clir.AddFlags(&flags)
	c.clir.Action(func() error {
		switch {
		case flags == emptyFlags || flags.List:
			return templates()
		default:
			targets := parseTargets(flags.Targets)
			return create(c.ctx, flags.Dir, targets)
		}
	})
	return
}

func parseTargets(targets string) (out map[string]string) {
	out = make(map[string]string)
	for _, s := range strings.Split(targets, " ") {
		chunks := strings.Split(s, ":")
		template := chunks[0]
		name := template
		if len(chunks) > 1 {
			name = chunks[1]
		}
		out[template] = name
	}
	return
}
