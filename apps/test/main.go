package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	factoryApphost "github.com/cryptopunkscc/portal/factory/apphost"
	"github.com/cryptopunkscc/portal/runner/cli"
	rpcApphost "github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"log"
	"os"
)

func main() {
	log.Println("main test args", os.Args)

	cli.Run(cmd.Handler{
		Name: "test",
		Desc: "Portal test app",
		Func: func() error { return nil },
		Sub: cmd.Handlers{
			{
				Name: "echo e",
				Func: echo,
			},
			{
				Name: "s",
				Desc: "serve",
				Func: func() (s string, err error) {
					log.Println("test serve")
					err = rpcApphost.Rpc(factoryApphost.Default()).Router(
						cmd.Handler{
							Func: echo,
						},
						apphost.NewPort("test"),
					).Run(context.Background())
					return "serve_end", err
				},
			},
		},
	})

	//time.Sleep(10 * time.Second)
}

func echo(args ...string) []string {
	return append(args, "echo123")
}
