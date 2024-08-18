package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Install(src string) (err error) {
	file, err := source.File(src)
	if err != nil {
		return err
	}
	return InstallSource(file)
}

func InstallSource(source target.Source) (err error) {
	for _, bundle := range apps.Resolver[target.Bundle_]().List(source) {
		if err = install(bundle); err != nil {
			log.Printf("Error copying file %s: %v", bundle.Abs(), err)
		}
	}
	return
}

func install(bundle target.Bundle_) error {
	pkg := bundle.Package()
	name := filepath.Base(bundle.Abs())
	dstPath := filepath.Join(portalAppsDir, name)
	println(fmt.Sprintf("* copying %s to %s [DONE]", name, dstPath))
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
