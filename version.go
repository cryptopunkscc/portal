package portal

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/git"
)

func GoModuleVersion() (version string) {
	version = Version
	if hash, err := git.TimestampHash(); err == nil {
		version += "-" + hash
	}
	return
}
