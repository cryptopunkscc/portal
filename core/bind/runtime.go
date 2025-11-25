package bind

import (
	"context"

	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/target"
)

type NewCore func(ctx context.Context, portal target.Portal_) (Core, context.Context)

type Core interface {
	Apphost
	bind.Process
}
