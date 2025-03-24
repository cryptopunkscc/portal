package core

import (
	. "github.com/cryptopunkscc/portal/api/mobile"
)

func (m *service) Start() {
	m.mobile.Event(&Event{Msg: STARTING})
	if err := m.Service.Start(m.ctx); err != nil {
		m.mobile.Event(&Event{Msg: STOPPED, Err: err})
		return
	}
	m.mobile.Event(&Event{Msg: STARTED})
	go func() {
		err := m.Wait()
		m.mobile.Event(&Event{Msg: STOPPED, Err: err})
	}()
	return
}
