package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"io/fs"
)

func RawTargets(files fs.FS) <-chan PortalRawModule {
	return Find[PortalRawModule](files, ".")
}

func DevTargets(files fs.FS) <-chan PortalNodeModule {
	return Find[PortalNodeModule](files, ".")
}

func BundleTargets(files fs.FS, dir string) <-chan Bundle {
	return Find[Bundle](files, dir)
}

func ProdTargets(files fs.FS) <-chan runner.Target {
	return FindAll(files, ".", PortalRawModule{}, Bundle{})
}
