package core

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/core/portald"
)

type service struct {
	portald.Service
	ctx    context.Context
	mobile Api
}

var _ Core = &service{}
