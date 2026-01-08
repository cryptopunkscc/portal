package main

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/runner/v2/goja"
	goja_dist "github.com/cryptopunkscc/portal/runner/v2/goja/dist"
	goja_pro "github.com/cryptopunkscc/portal/runner/v2/goja/pro"
	"github.com/cryptopunkscc/portal/runner/v2/wails"
	wails_dist "github.com/cryptopunkscc/portal/runner/v2/wails/dist"
	wails_pro "github.com/cryptopunkscc/portal/runner/v2/wails/pro"
	"github.com/cryptopunkscc/portal/source"
)

func runTarget(ctx context.Context, src string, args ...string) (err error) {
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
