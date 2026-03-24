package core

import (
	"errors"

	"github.com/cryptopunkscc/portal/mobile"
	"github.com/cryptopunkscc/portal/pkg/bind"
	bind2 "github.com/cryptopunkscc/portal/pkg/bind/src"
	"github.com/cryptopunkscc/portal/pkg/source/app"
	"github.com/spf13/afero"
)

func (srv *Service) App(pkg string) (mobile.App, error) {
	bundle, err := srv.app(pkg)
	if err != nil {
		return nil, err
	}
	return &App{
		Bundle: bundle,
		core:   bind2.DefaultCoreFactory{}.Create(srv.ctx),
		assets: assets{afero.IOFS{Fs: bundle.Fs}},
		manifest: &mobile.Manifest{
			Name:        bundle.Name,
			Title:       bundle.Title,
			Description: bundle.Description,
			Package:     bundle.Package,
			Icon:        bundle.Icon,
		},
	}, nil
}

func (srv *Service) app(src string) (bundle *app.Bundle, err error) {
	bundle, err = app.Objects{Astrald: &srv.client}.GetAppBundle(srv.ctx, src)
	err = errors.Unwrap(err)
	return
}

type App struct {
	*app.Bundle
	core     bind.Core
	assets   mobile.Assets
	manifest *mobile.Manifest
}

func (a App) Manifest() *mobile.Manifest {
	return a.manifest
}

func (a App) Assets() mobile.Assets {
	return a.assets
}

func (a App) Core() bind.Core {
	return a.core
}
