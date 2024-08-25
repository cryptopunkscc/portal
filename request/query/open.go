package query

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

type Opener struct{ target.Port }

var Open = Opener{target.PortOpen}

func (port Opener) Run(ctx context.Context, app string) (packages []string, err error) {
	plog.Get(ctx).Type(port).Println("Running", app)
	session, err := rpc.NewFlow().Query(id.Anyone, port.Base)
	if err != nil {
		err = fmt.Errorf("cannot query %s: %w", port.Base, err)
		return
	}
	defer session.Close()
	if packages, err = rpc.Query[[]string](session, port.Name, app); err != nil {
		err = fmt.Errorf("cannot open %s: %w", app, err)
		return
	}
	select {
	case <-ctx.Done():
	case <-rpc.Done(session):
	}
	return
}

func (port Opener) Start(ctx context.Context, app string) (packages []string, err error) {
	plog.Get(ctx).Type(port).Println("Starting", app)
	request := rpc.NewRequest(id.Anyone)
	if packages, err = rpc.Query[[]string](request, port.String(), app); err != nil {
		err = fmt.Errorf("cannot open %s: %w", app, err)
		return
	}
	return
}
