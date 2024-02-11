package apphost

import "embed"

//go:embed js/*
var jsFs embed.FS

func GenerateAppHostJs(modules ...string) (ah string, err error) {
	for _, module := range modules {
		var buff []byte
		if buff, err = jsFs.ReadFile("js/" + module + ".js"); err != nil {
			return
		}
		ah = ah + string(buff) + "\n"
	}
	return
}
