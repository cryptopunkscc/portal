package install

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
)

func Runner(appsDir string) Install {
	return Install{appsDir: appsDir}
}

type Install struct {
	appsDir string
}

func (i Install) Run(src string) (c <-chan Result, err error) {
	file, err := source.File(src)
	if err != nil {
		return
	}
	results := make(chan Result)
	c = results
	go i.All(file, results)
	return
}

type Result struct {
	Manifest target.Manifest
	Error    error
}

func (r Result) MarshalCLI() string {
	status := "SUCCESS"
	if r.Error != nil {
		status = "FAILURE: " + r.Error.Error()
	}
	return fmt.Sprintf("install %s: %s", r.Manifest.Name, status)
}
