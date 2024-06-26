package version

import (
	"github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/pkg/git"
)

func GoModuleVersion() (version string) {
	version = portal.Version
	if hash, err := git.TimestampHash(); err == nil {
		version += "-" + hash
	}
	return
}
