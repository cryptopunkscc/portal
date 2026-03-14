package portal_goja

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/bind/src"
	"github.com/cryptopunkscc/portal/pkg/client"
	"github.com/cryptopunkscc/portal/pkg/runner/goja"
	"github.com/cryptopunkscc/portal/pkg/source"
	"github.com/cryptopunkscc/portal/pkg/source/app"
)

type Application struct {
	Astrald *client.Astrald
}

func (a Application) Run(ctx context.Context, src string, args ...string) (err error) {
	s := source.Providers{
		source.OsFs,
		app.Objects{Astrald: a.Astrald},
	}.GetSource(src)
	if s == nil {
		return fs.ErrNotExist
	}

	f := bind.DefaultCoreFactory{Astrald: a.Astrald}
	for _, ss := range source.Collect(s,
		&goja.AppRunner{},
		&goja.BundleRunner{},
	) {
		ctx := f.Create(ctx)
		return ss.(goja.Runner).Run(ctx, args...)
	}

	return fs.ErrInvalid
}
