package install

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
)

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
	Id       int
	Manifest target.Manifest
	Error    error
}

func (r Result) MarshalCLI() string {
	status := "[DONE]"
	if r.Error != nil {
		status = "[FAILURE]: " + r.Error.Error()
	}
	return fmt.Sprintf("%d. %s %s\n", r.Id, r.Manifest.Name, status)
}
