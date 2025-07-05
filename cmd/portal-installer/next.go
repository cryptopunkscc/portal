package main

func nextInstallation() (err error) {
	if err = installBinaries(); err != nil {
		return
	}
	if err = portalRun(); err != nil {
		return
	}
	if err = installApps(); err != nil {
		return
	}
	if err = portalRun("close"); err != nil {
		return
	}
	return
}
