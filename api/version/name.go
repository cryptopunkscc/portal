package version

import (
	_ "embed"
	"github.com/cryptopunkscc/portal/api"
	"github.com/cryptopunkscc/portal/pkg/git"
	"github.com/cryptopunkscc/portal/pkg/vcs"
	"strings"
)

//go:embed name
var version string

func Name() string {
	if version = strings.TrimSpace(version); version == "" {
		version = Resolve()
	}
	return version
}

func Resolve() (version string) {
	version = goModuleVersion()
	if vcs.ReadBuildInfo().Modified != "" {
		version += " [MODIFIED]"
	}
	return
}

func goModuleVersion() (version string) {
	version = api.Version
	if hash, err := git.TimestampHash(); err == nil {
		version += "-" + hash
	}
	return
}
