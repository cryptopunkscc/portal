package main

import "os"

func (d *Install) buildInstaller() {
	mkdirBin()
	buildAstral()
	buildAnc()
	buildPortal()
	buildPortalDev()
	buildPortalInstaller()
}

func mkdirBin() {
	if err := os.Mkdir("./cmd/portal-installer/bin/", 0755); err != nil {
		panic(err)
	}
}
