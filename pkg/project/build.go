package project

import (
	"os"
)

func BuildPortalApps(root, dir string, dependencies ...NodeModule) (err error) {
	for m := range Find[PortalNodeModule](os.DirFS(root), dir) {
		if !m.CanNpmRunBuild() {
			continue
		}
		if err = m.PrepareBuild(dependencies...); err != nil {
			return err
		}
	}
	return
}
