package install

func Runner(appsDir string) Install {
	return Install{appsDir: appsDir}
}

type Install struct {
	appsDir string
}
