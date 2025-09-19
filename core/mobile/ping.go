package core

import . "github.com/cryptopunkscc/portal/api/mobile"

func (m *service) Ping() {
	switch {
	case m.status == FRESH && m.HasUser():
		m.set(STARTED)
	default:
		m.set(m.status)
	}
}
