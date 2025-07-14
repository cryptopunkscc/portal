package apps

import (
	"context"
	"github.com/cryptopunkscc/portal/core/portald/test/apps/go_client"
	"github.com/cryptopunkscc/portal/core/portald/test/apps/go_service"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"testing"
	"time"
)

func TestGoRpc(t *testing.T) {
	t.SkipNow() // FIXME

	ctx := context.Background()
	log := plog.New().Set(&ctx)
	srv := go_service.NewTestGoService("test.go")
	start := rpc.Start(srv.Router)

	if err := start(ctx); err != nil {
		log.P().Println(err)
	}

	time.Sleep(100 * time.Millisecond)

	go_client.NewTestClient("test", "go").Run(t)
}
