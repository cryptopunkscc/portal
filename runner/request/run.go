package request

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Open(ctx context.Context, src string, _ ...string) (err error) {
	log := plog.Get(ctx)
	log.Println("starting query", "portald.open", src)
	request := apphost.Default.Rpc().Request("portald")
	err = rpc.Command(request, "open", src)
	if err != nil {
		log.E().Printf("cannot query %s: %v", src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	log.Println("started query", "portald.open", src)
	return
}
