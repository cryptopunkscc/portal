package install

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"io"
	"os"
	"path/filepath"
)

func (i Runner) BundlesByPath(src string) (c <-chan Result, err error) {
	file, err := source.File(src)
	if err != nil {
		return
	}
	results := make(chan Result)
	c = results
	go i.installBundles(file, results)
	return
}

func (i Runner) installBundles(source target.Source, c chan<- Result) {
	defer close(c)
	for id, bundle := range apps.Resolver[target.Bundle_]().List(source) {
		err := i.Bundle(bundle)
		c <- Result{
			Id:       id,
			Error:    err,
			Manifest: *bundle.Manifest(),
		}
	}
	return
}

func (i Runner) Bundle(bundle target.Bundle_) error {
	if _, err := (token.Repository{}).Resolve(bundle.Manifest().Package); err != nil {
		return err
	}
	pkg := bundle.Package()
	name := filepath.Base(bundle.Abs())
	dstPath := filepath.Join(i.OutputDir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	src, err := pkg.Files().Open(pkg.Path())
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(dst, src)
	return err
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
