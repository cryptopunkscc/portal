package main

import (
	"context"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
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
					router := apphost.Default.Rpc().Router(
						cmd.Handler{
							Func: echo,
						},
					)
					if err = router.Init(context.Background()); err != nil {
						return
					}
					err = router.Listen()
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
