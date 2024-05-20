package clir

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
	"log"
	"os"
	"reflect"
	"strings"
)

func (c Cli) Apps() {
	flags := struct {
		Path string `pos:"1" default:"."`
	}{}
	cmd := c.clir.NewSubCommand("t", "Print all targets in given directory.")
	cmd.AddFlags(&flags)
	cmd.Action(func() (err error) {
		wd, _ := os.Getwd()
		for source := range sources.FromPath[target.Source](flags.Path) {
			log.Println(reflect.TypeOf(source), "\t", strings.TrimPrefix(source.Abs(), wd+"/"))
		}
		return
	})
	return
}
