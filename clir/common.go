package clir

import (
	"github.com/leaanthony/clir"
	"log"
)

type Version func() string

func (c cli) Version(
	handle Version,
) (cmd *clir.Command) {
	cmd = c.clir.NewSubCommand("v", "Print portal version.")
	cmd.Action(func() (_ error) {
		log.Println(handle())
		return
	})
	return
}
