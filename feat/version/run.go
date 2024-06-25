package version

import (
	"github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/pkg/vcs"
)

func Run() (version string) {
	version = portal.GoModuleVersion()
	if vcs.ReadBuildInfo().Modified != "" {
		version += " [MODIFIED]"
	}
	return
}
