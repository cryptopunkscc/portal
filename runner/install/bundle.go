package install

import (
	"github.com/cryptopunkscc/portal/api/target"
	"io"
	"os"
	"path/filepath"
)

func (i Install) Bundle(bundle target.Bundle_) error {
	if _, err := Token(bundle.Manifest().Package); err != nil {
		return err
	}
	pkg := bundle.Package()
	name := filepath.Base(bundle.Abs())
	dstPath := filepath.Join(i.appsDir, name)
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
