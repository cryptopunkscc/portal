package main

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test/rpc"
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
