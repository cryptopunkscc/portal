package main

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/bind/src"
	"github.com/cryptopunkscc/portal/pkg/runner/goja"
	"github.com/cryptopunkscc/portal/pkg/runner/goja/dist"
	"github.com/cryptopunkscc/portal/pkg/runner/goja/pro"
	"github.com/cryptopunkscc/portal/pkg/runner/wails"
	"github.com/cryptopunkscc/portal/pkg/runner/wails/dist"
	"github.com/cryptopunkscc/portal/pkg/runner/wails/pro"
	"github.com/cryptopunkscc/portal/pkg/source"
	"github.com/cryptopunkscc/portal/pkg/util/os"
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
			return r.Run(ctx, args...)
		case wails.Runner:
			ctx := bind.DefaultCoreFactory{}.Create(ctx)
			return r.Run(&Adapter{ctx})
		}
	}

	return fs.ErrInvalid
}

type Adapter struct{ *bind.Core }
