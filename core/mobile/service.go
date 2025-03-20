package core

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/portald"
)

type service struct {
	portald.Service[target.App_]
	mobile  Api
	astrald astrald
	Ctx     context.Context
	Find    target.Find[target.App_]
}

var _ Portald = &service{}

func (m *service) Stop() {
	m.Shutdown()
	m.Ctx = context.Background()
}
