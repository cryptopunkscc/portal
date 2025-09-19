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
	status int32
}

var _ Core = &service{}

func (m *service) set(status int32) {
	m.status = status
	m.mobile.Status(status)
}

func (m *service) err(err error) {
	if err != nil {
		m.mobile.Error(err.Error())
	}
}
