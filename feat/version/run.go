package version

import (
	_ "embed"
	"github.com/cryptopunkscc/portal/pkg/vcs"
	"strings"
)

//go:embed name
var version string

func Run() string {
	if version = strings.TrimSpace(version); version == "" {
		version = Resolve()
	}
	return version
}

func Resolve() (version string) {
	version = GoModuleVersion()
	if vcs.ReadBuildInfo().Modified != "" {
		version += " [MODIFIED]"
	}
	return
}
