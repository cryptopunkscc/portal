package main

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/os"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja/dist"
	"github.com/cryptopunkscc/portal/runner/goja/pro"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails/dist"
	"github.com/cryptopunkscc/portal/runner/wails/pro"
	"github.com/cryptopunkscc/portal/source"
)

func runTarget(ctx context.Context, src string, args ...string) (err error) {
	src = os.Abs(src)
	s := source.Providers{
		source.OsFs,
	}.GetSource(src)
	if s == nil {
		return fs.ErrNotExist
	}

	for _, ss := range source.Collect(s,
		&goja_pro.Runner{},
		&goja_dist.Runner{},
		&goja.BundleRunner{},
		&wails_pro.Runner{},
		&wails_dist.Runner{},
		&wails.BundleRunner{},
	) {
		switch r := ss.(type) {
		case goja.Runner:
			ctx := bind.DefaultCoreFactory{}.Create(ctx)
			return r.Run(*ctx, args...)
		case wails.Runner:
			ctx := bind.DefaultCoreFactory{}.Create(ctx)
			return r.Run(&Adapter{ctx})
		}
	}

	return fs.ErrInvalid
}

type Adapter struct{ bind.Core }
