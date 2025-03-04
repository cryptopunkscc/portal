package request

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc"
	"github.com/cryptopunkscc/portal/runtime/rpc/apphost"
)

func Open(ctx context.Context, src string, _ ...string) (err error) {
	log := plog.Get(ctx)
	log.Println("starting query", "portal.open", src)
	request := apphost.Request("portal")
	err = rpc.Command(request, "open", src)
	if err != nil {
		log.E().Printf("cannot query %s: %v", src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	log.Println("started query", "portal.open", src)
	return
}
