package rpc

import (
	"context"
	jrpc "github.com/cryptopunkscc/go-apphost-jrpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"log"
	"os/exec"
)

func Run(bindings runner.Bindings) (err error) {

	s := jrpc.NewApp("portal")
	s.Logger(log.New(log.Writer(), "service ", 0))
	s.With(bindings)
	s.RouteFunc("open", open)

	ctx := context.Background()
	if err = s.Run(ctx); err != nil {
		return
	}
	<-ctx.Done()
	return
}

func open(path string, background bool) (pid int, err error) {
	c := exec.Command("portal", path)
	if !background {
		err = c.Run()
		return
	}
	if err = c.Start(); err != nil {
		return
	}
	pid = c.Process.Pid
	return
}
