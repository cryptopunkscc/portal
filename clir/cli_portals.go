package clir

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"log"
	"os"
	"reflect"
	"strings"
)

func (c Cli) Portals(
	find target.Find[target.Portal],
) {
	flags := struct {
		Path string `pos:"1" default:"."`
	}{}
	cmd := c.clir.NewSubCommand("t", "Print all targets in given directory.")
	cmd.AddFlags(&flags)
	cmd.Action(func() (err error) {
		wd, _ := os.Getwd()
		portals, err := find(c.ctx, flags.Path)
		if err != nil {
			return
		}
		for _, source := range portals {
			log.Println(reflect.TypeOf(source), "\t", strings.TrimPrefix(source.Abs(), wd+"/"))
		}
		return
	})
	return
}