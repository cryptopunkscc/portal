package core

import (
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
)

func (m *service) Install() (err error) {
	for _, dist := range m.Resolve.List(
		source.Embed(apps.LauncherSvelteFS),
	) {
		if err = m.Service.Install().CopyOf(dist); err != nil {
			return
		}
	}
	return
}
