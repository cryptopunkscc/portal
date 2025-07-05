package main

func nextInstallation() (err error) {
	if err = installBinaries(); err != nil {
		return
	}
	return
}
