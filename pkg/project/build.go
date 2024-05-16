package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"path"
)

func BuildPortalApps(root, dir string, dependencies ...target.NodeModule) (err error) {
	for m := range FindInPath[*PortalNodeModule](path.Join(root, dir)) {

		if !m.CanNpmRunBuild() {
			continue
		}
		if err = m.PrepareBuild(dependencies...); err != nil {
			return err
		}
	}
	return
}
