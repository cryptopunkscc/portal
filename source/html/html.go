package html

import (
	"path"

	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/cryptopunkscc/portal/source/npm"
	"github.com/spf13/afero"
)

type App struct {
	app.Dist
	Html
}

func (a App) New() (src source.Source) {
	return &a
}

func (a *App) ReadSrc(src source.Source) (err error) {
	return source.Readers{&a.Dist, &a.Html}.ReadSrc(src)
}

func (a *App) WriteRef(ref source.Ref) (err error) {
	return source.Writers{&a.Dist, &a.Html}.WriteRef(ref)
}

type Project struct {
	npm.Project
	Html
	App App
}

func (p Project) New() (src source.Source) {
	return &p
}

func (p *Project) ReadSrc(src source.Source) (err error) {
	return source.Readers{&p.Html, &p.Project}.ReadSrc(src)
}

func (p *Project) WriteRef(ref source.Ref) (err error) {
	return source.Writers{&p.Project, &p.Html}.WriteRef(ref)
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
	return b.App.ReadSrc(&source.Ref{Fs: b.Zip.Unpacked})
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

type Html struct {
	IndexHtml source.Blob `json:"index,omitempty" yaml:"index,omitempty"`
}

func (h *Html) ReadSrc(src source.Source) (err error) {
	ref := *src.Ref_()
	ref.Path = path.Join(ref.Path, "index.html")
	return h.IndexHtml.ReadSrc(&ref)
}

func (h *Html) WriteRef(ref source.Ref) (err error) {
	ref.Path = path.Join(ref.Path, "index.html")
	return h.IndexHtml.WriteRef(ref)
}
