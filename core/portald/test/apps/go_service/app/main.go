package main

import (
	"context"
	"github.com/cryptopunkscc/portal/core/portald/test/apps/go_service"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func main() {
	ctx := context.Background()
	log := plog.New().Set(&ctx)
	srv := go_service.NewTestGoService("test.go")
	log.Type(srv).Println("start test.go")
	err := srv.Run(ctx)
	if err != nil {
		log.P().Println(err)
	}
}
