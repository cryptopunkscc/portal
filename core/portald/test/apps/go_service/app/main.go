package main

import (
	"context"

	"github.com/cryptopunkscc/portal/core/portald/test/apps/go_service"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func main() {
	ctx := context.Background()
	log := plog.New().Set(&ctx)
	srv := go_service.NewTestGoService("test.go")
	run := rpc.Run(srv)
	log.Type(srv).Println("start test.go")
	err := run(ctx)
	if err != nil {
		log.P().Println(err)
	}
}
