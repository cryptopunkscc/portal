package clir

import (
	"context"
)

type Build func(context.Context, string) error
type Clean func(string) error

func (c Cli) Build(
	build Build,
	clean Clean,
) {
	flags := struct {
		Path  string `pos:"1" default:"."`
		Clean bool   `name:"c" description:"Clean target directories from build artifacts without building."`
	}{}
	c.clir.AddFlags(&flags)
	c.clir.Action(func() (err error) {
		if flags.Clean {
			return clean(flags.Path)
		} else {
			_ = clean(flags.Path)
			return build(c.ctx, flags.Path)
		}
	})
	return
}
