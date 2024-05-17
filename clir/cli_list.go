package clir

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"log"
)

type List func() []target.App

func (c Cli) List(handle List) {
	cmd := c.clir.NewSubCommand("l", "List installed apps.")
	cmd.Action(func() (_ error) {
		for i, app := range handle() {
			log.Println(i, app.Manifest())
		}
		return
	})
	return
}
