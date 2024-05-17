package clir

import (
	"context"
	"github.com/leaanthony/clir"
	"log"
)

type Cli struct {
	clir    *clir.Cli
	ctx     context.Context
	version Version
}

func NewCli(
	ctx context.Context,
	name, description string,
	version Version,
) *Cli {
	c := clir.NewCli(name, description, version())
	c.Version()
	return &Cli{
		ctx:     ctx,
		clir:    c,
		version: version,
	}
}

type Version func() string

func (c Cli) Version() (cmd *clir.Command) {
	cmd = c.clir.NewSubCommand("v", "Print portal version.")
	cmd.Action(func() (_ error) {
		log.Println(c.version())
		return
	})
	return
}

func (c Cli) Run() error {
	return c.clir.Run()
}
