package core

import (
	"context"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
	. "github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (m *service) App(path string) mobile.App {
	apps, err := m.find(m.ctx, path)
	if err != nil {
		plog.Get(m.ctx).Type(m).E().Println(err)
		return nil
	}
	if len(apps) == 0 {
		plog.Get(m.ctx).Type(m).E().Println(target.ErrNotFound)
		return nil
	}
	source := apps[0]
	return &app_{
		ctx:     m.ctx,
		source:  source,
		newCore: m.cores().NewFrontendFunc(),
	}
}

func (m *service) cores() CoreFactory {
	return CoreFactory{Repository: *m.Tokens()}
}

type app_ struct {
	ctx     context.Context
	source  target.App_
	newCore NewCore
	core    Core
}

var _ mobile.App = &app_{}

func (a *app_) Assets() mobile.Assets {
	return assets{a.source.FS()}
}

func (a *app_) Core() bind.Core {
	if a.core == nil {
		a.core, a.ctx = a.newCore(a.ctx, a.source)
	}
	return a.core
}

func (a *app_) Manifest() *mobile.Manifest {
	manifest := a.source.Manifest()
	return &mobile.Manifest{
		Name:        manifest.Name,
		Title:       manifest.Title,
		Description: manifest.Description,
		Package:     manifest.Package,
		Version:     manifest.Version,
		Icon:        manifest.Icon,
	}
}
