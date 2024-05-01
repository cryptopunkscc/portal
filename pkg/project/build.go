package project

import (
	"os"
	"path"
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

func Path(src string) (base, sub string, err error) {
	src = path.Clean(src)
	base = src
	sub = "."
	if !path.IsAbs(base) {
		sub = src
		base, err = os.Getwd()
		if err != nil {
			return
		}
	}
	return
}
