package core

import (
	"context"

	_ "github.com/cryptopunkscc/astrald/mod/all"
	. "github.com/cryptopunkscc/portal/mobile"
)

type service struct {
	ctx    context.Context
	mobile Api
	status int32
}

func (m *service) Stop() { /*TODO*/ }

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
