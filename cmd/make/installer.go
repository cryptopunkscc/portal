package main

func (d *Install) buildInstaller() {
	buildAstral()
	buildPortal()
	buildPortalDev()
	buildPortalInstaller()
}
