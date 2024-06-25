package portal

import (
	"github.com/cryptopunkscc/portal/pkg/git"
)

func GoModuleVersion() (version string) {
	version = Version
	if hash, err := git.TimestampHash(); err == nil {
		version += "-" + hash
	}
	return
}
