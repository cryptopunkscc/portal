package install

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"io"
	"os"
	"path/filepath"
)

func Runner(appsDir string) Install {
	return Install{appsDir: appsDir}
}

type Install struct {
	appsDir string
}

func (i Install) Run(ctx context.Context, src string, _ ...string) error {
	file, err := source.File(src)
	if err != nil {
		return err
	}
	return i.All(ctx, file)
}

func (i Install) All(ctx context.Context, source target.Source, _ ...string) error {
	log := plog.Get(ctx).Type(i)
	for _, bundle := range apps.Resolver[target.Bundle_]().List(source) {
		if err := i.Bundle(ctx, bundle); err != nil {
			log.Printf("Error copying file %s: %v", bundle.Abs(), err)
		}
	}
	return nil
}

func (i Install) Bundle(_ context.Context, bundle target.Bundle_, _ ...string) error {
	pkg := bundle.Package()
	name := filepath.Base(bundle.Abs())
	dstPath := filepath.Join(i.appsDir, name)
	_, _ = fmt.Fprintf(os.Stdout, "* copying %s to %s [DONE]\n", name, dstPath)
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
