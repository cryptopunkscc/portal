package main

import (
	"errors"
	"io/fs"
	"os"
)

func (d *Install) buildInstaller() {
	mkdirBin()
	buildAstral()
	buildAnc()
	buildPortal()
	buildPortalDev()
	buildPortalInstaller()
}

func mkdirBin() {
	err := os.Mkdir("./cmd/portal-installer/bin/", 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic(err)
	}
}
