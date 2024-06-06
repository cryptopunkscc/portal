package rpc

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"testing"
	"time"
)

func TestGoRpc(t *testing.T) {
	ctx := context.Background()
	log := plog.New().Set(&ctx)
	srv := NewTestGoService("test.go")

	if err := srv.Start(ctx); err != nil {
		log.P().Println()
	}

	time.Sleep(100 * time.Millisecond)

	NewTestClient("test.go").Run(t)
}
