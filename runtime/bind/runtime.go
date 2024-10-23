package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/target"
)

type NewRuntime func(ctx context.Context, portal target.Portal_) Runtime

type Runtime interface {
	Apphost
	bind.Sys
}

type Module struct {
	Apphost
	bind.Sys
}
