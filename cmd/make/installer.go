package main

func (d *Install) buildInstaller() {
	buildAstral()
	buildAnc()
	buildPortal()
	buildPortalDev()
	buildPortalInstaller()
}
