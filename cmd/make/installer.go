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
	goPortal.Build()
	goPortalApp.Build()
	goPortalTray.Build()
	goPortalAppGoja.Build()
	goPortalAppWails.Build()
	goPortalDev.Build()
	goPortalDevGo.Build()
	goPortalDevGoja.Build()
	goPortalDevExec.Build()
	goPortalDevWails.Build()
	buildPortalInstaller()
}

func mkdirBin() {
	err := os.Mkdir("./cmd/portal-installer/bin/", 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic(err)
	}
}
