package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/integration_test/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
)

func main() {
	ctx := context.Background()
	log := plog.New().Set(&ctx)
	srv := rpc.NewTestGoService("test.go.service")
	log.Type(srv).Println("start")
	err := srv.Run(ctx)
	if err != nil {
		log.P().Println()
	}
}
