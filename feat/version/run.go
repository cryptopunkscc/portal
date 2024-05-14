package version

import (
	portal "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/vcs"
)

func Run() (version string) {
	version = portal.GoModuleVersion()
	if vcs.ReadBuildInfo().Modified != "" {
		version += " [MODIFIED]"
	}
	return
}
