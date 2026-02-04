package core

import (
	"context"
	"errors"

	. "github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/mobile"
	"github.com/cryptopunkscc/portal/pkg/bind"
	bind2 "github.com/cryptopunkscc/portal/pkg/bind/src"
	"github.com/cryptopunkscc/portal/pkg/runner/goja"
	"github.com/cryptopunkscc/portal/pkg/source/app"
)

var ErrNotExist = errors.New("no apps found")

func (m *service) App(src string) (mobile.App, error) {
	s := app.Objects{Adapter: m.coreFactory.Adapter}.GetSource(src)
	if s == nil {
		return nil, ErrNotExist
	}
	b := &goja.BundleRunner{}
	if err := b.ReadSrc(s); err != nil {
		return nil, err
	}
	return &app_{
		ctx:         m.ctx,
		source:      b,
		coreFactory: m.coreFactory,
	}, nil
}

func (m *service) newCore(ctx context.Context, portal app.App) (*bind2.Core, context.Context) {
	var c = AutoTokenCoreFactory{
		PkgName: portal.GetMetadata().Package,
		Tokens:  m.Tokens(),
	}.Create(ctx)
	return c, c.Context
}

type app_ struct {
	ctx         context.Context
	source      app.App
	coreFactory bind2.DefaultCoreFactory
	core        *bind2.Core
}

var _ mobile.App = &app_{}

func (a *app_) Assets() mobile.Assets {
	return assets{a.source.Ref_().PathFS()}
}

func (a *app_) Core() bind.Core {
	if a.core == nil {
		a.core = a.coreFactory.Create(a.ctx)
		a.ctx = a.core
	}
	return a.core
}

func (a *app_) Manifest() *mobile.Manifest {
	manifest := a.source.GetMetadata()
	return &mobile.Manifest{
		Name:        manifest.Name,
		Title:       manifest.Title,
		Description: manifest.Description,
		Package:     manifest.Package,
		Version:     manifest.Version(),
		Icon:        manifest.Icon,
	}
}
