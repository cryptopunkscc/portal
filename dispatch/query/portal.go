package query

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
)

type Portal struct{ port target.Port }

func NewPortal() Portal {
	return Portal{port: target.PortPortal}
}

func (p Portal) Start(ctx context.Context, src string, args ...string) (err error) {
	request := rpc.NewRequest(id.Anyone, p.port.String())
	request.Logger(plog.Get(ctx).Type(p))
	anyArgs := make([]any, len(args))
	for i, arg := range args {
		anyArgs[i] = arg
	}
	return rpc.Command(request, src, anyArgs...)
}
