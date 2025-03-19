package main

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/core/rpc/cli"
	"log"
)

func (a Application) Arg(ctx context.Context, runner *cli.Runner, arg string, cmd ...string) (err error) {
	if arg == "" {
		return errors.New("empty arg")
	}
	for _, s := range cmd {
		r := *runner
		c := s + " " + arg
		log.Println("Running command: " + c)
		r.Setup(c)
		r.Dependencies = append([]any{ctx}, r.Dependencies...)
		if err = r.Respond(&r.Conn); err != nil {
			return
		}
	}
	return
}
