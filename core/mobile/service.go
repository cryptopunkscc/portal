package core

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/portald"
)

type service struct {
	portald.Service[target.App_]
	ctx    context.Context
	mobile Api
}

var _ Core = &service{}
