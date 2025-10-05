package core

import (
	. "github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (m *service) Start() {
	m.set(STARTING)

	if err := m.Service.Start(m.ctx); err != nil {
		plog.Println(err)
		m.err(err)
		m.set(STOPPED)
		return
	}

	_ = m.installApps()
	m.set(STARTED)

	go func() {
		err := m.Wait()
		plog.Println(err)
		m.err(err)
		m.set(STOPPED)
	}()
	return
}
