package main

func firstInstallation(username string) (err error) {
	if err = installBinaries(); err != nil {
		return
	}
	if err = portalRun(); err != nil {
		return
	}
	if err = portalRun("user", "create", username); err != nil {
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
