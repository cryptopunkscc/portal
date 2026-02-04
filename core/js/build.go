package js

import (
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/npm"
)

func BuildPortalLib() (err error) {
	dir, err := golang.FindProjectRoot()
	if err != nil {
		return
	}
	nm := npm.NodeModule{}
	if err = nm.ReadSrc(source.OSRef(dir, "core", "js")); err != nil {
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
