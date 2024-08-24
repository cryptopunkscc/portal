package clir

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"log"
)

type List func(ctx context.Context) (target.Portals[target.App_], error)

func (c Cli) List(handle List) {
	cmd := c.clir.NewSubCommand("l", "List installed apps.")
	cmd.Action(func() (err error) {
		list, err := handle(c.ctx)
		if err != nil {
			return
		}
		for i, app := range list {
			log.Println(i, app.Manifest())
		}
		return
	})
	return
}
