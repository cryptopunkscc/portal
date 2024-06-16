package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/test/rpc"
)

func main() {
	ctx := context.Background()
	log := plog.New().Set(&ctx)
	srv := rpc.NewTestGoService("test.go")
	log.Type(srv).Println("start test.go")
	err := srv.Run(ctx)
	if err != nil {
		log.P().Println(err)
	}
}
