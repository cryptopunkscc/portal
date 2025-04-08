package rpc

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"testing"
	"time"
)

func TestGoRpc(t *testing.T) {
	ctx := context.Background()
	log := plog.New().Set(&ctx)
	srv := NewTestGoService("test.go")

	if err := srv.Router.Start(ctx); err != nil {
		log.P().Println(err)
	}

	time.Sleep(100 * time.Millisecond)

	NewTestClient("test", "go").Run(t)
}
