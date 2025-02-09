package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/target"
)

type NewRuntime func(ctx context.Context, portal target.Portal_) (Runtime, context.Context)

type Runtime interface {
	Apphost
	bind.Sys
}

type Module struct {
	Apphost
	bind.Sys
}
