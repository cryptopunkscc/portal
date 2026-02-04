package js

import (
	"github.com/cryptopunkscc/portal/pkg/source"
	"github.com/cryptopunkscc/portal/pkg/source/npm"
	"github.com/cryptopunkscc/portal/pkg/util/go"
)

func BuildPortalLib() (err error) {
	dir, err := golang.FindProjectRoot()
	if err != nil {
		return
	}
	nm := npm.NodeModule{}
	if err = nm.ReadSrc(source.OSRef(dir, "pkg", "bind", "js")); err != nil {
		return
	}
	if err = nm.NpmInstall(); err != nil {
		return
	}
	if err = nm.Build(); err != nil {
		return
	}
	return
}
