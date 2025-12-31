package js

import (
	"path"

	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/cryptopunkscc/portal/source/npm"
	"github.com/spf13/afero"
)

type App struct {
	app.Dist
	Js
}

func (a App) New() (src source.Source) {
	return &a
}

func (a *App) ReadSrc(src source.Source) (err error) {
	return source.Readers{&a.Dist, &a.Js}.ReadSrc(src)
}

func (a *App) WriteRef(ref source.Ref) (err error) {
	return source.Writers{&a.Dist, &a.Js}.WriteRef(ref)
}

type Project struct {
	npm.Project
	Js  Js
	App App
}

var _ app.App = &Project{}

func (p Project) New() (src source.Source) {
	return &p
}

func (p Project) GetDist() (d app.Dist) {
	_ = d.ReadSrc(p.Sub("dist"))
	return
}

func (p *Project) ReadSrc(src source.Source) (err error) {
	return source.Readers{&p.Js, &p.Project}.ReadSrc(src)
}

func (p *Project) WriteRef(ref source.Ref) (err error) {
	return source.Writers{&p.Project, &p.Js}.WriteRef(ref)
}

type Bundle struct {
	App
	source.Zip
}

func (b Bundle) New() (src source.Source) {
	return &b
}

func (b *Bundle) ReadSrc(src source.Source) (err error) {
	if err = b.Zip.ReadSrc(src); err != nil {
		return
	}
	if err = b.App.ReadSrc(&source.Ref{Fs: b.Zip.Unpacked}); err != nil {
		return
	}
	return
}

func (b *Bundle) WriteRef(ref source.Ref) (err error) {
	if b.Fs == nil {
		b.Fs = afero.NewMemMapFs()
	}
	if err = b.App.WriteRef(b.Ref); err != nil {
		return
	}
	ref.Path = path.Join(ref.Path, b.BundleName())
	b.Zip.Unpacked = b.Fs
	return b.Zip.WriteRef(ref)
}

type Js struct {
	MainJs source.Blob `json:"main,omitempty" yaml:"main,omitempty"`
}

func (j *Js) ReadSrc(src source.Source) (err error) {
	ref := *src.Ref_()
	ref.Path = path.Join(ref.Path, "main.js")
	return j.MainJs.ReadSrc(&ref)
}

func (j *Js) WriteRef(ref source.Ref) (err error) {
	ref.Path = path.Join(ref.Path, "main.js")
	return j.MainJs.WriteRef(ref)
}
