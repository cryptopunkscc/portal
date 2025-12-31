package source

import (
	"path"

	"github.com/spf13/afero"
)

type HtmlProject struct {
	NpmProject
	Html
	App HtmlApp
}

func (p HtmlProject) New() (src Source) {
	return &p
}

func (p *HtmlProject) ReadSrc(src Source) (err error) {
	return Readers{&p.Html, &p.NpmProject}.ReadSrc(src)
}

func (p *HtmlProject) WriteRef(ref Ref) (err error) {
	return Writers{&p.NpmProject, &p.Html}.WriteRef(ref)
}

type HtmlApp struct {
	App
	Html
}

func (a HtmlApp) New() (src Source) {
	return &a
}

func (a *HtmlApp) ReadSrc(src Source) (err error) {
	return Readers{&a.App, &a.Html}.ReadSrc(src)
}

func (a *HtmlApp) WriteRef(ref Ref) (err error) {
	return Writers{&a.App, &a.Html}.WriteRef(ref)
}

type HtmlBundle struct {
	HtmlApp
	Zip
}

func (b HtmlBundle) New() (src Source) {
	return &b
}

func (d *HtmlBundle) ReadSrc(src Source) (err error) {
	if err = d.Zip.ReadSrc(src); err != nil {
		return
	}
	return d.HtmlApp.ReadSrc(&Ref{Fs: d.Zip.Unpacked})
}

func (d *HtmlBundle) WriteRef(ref Ref) (err error) {
	if d.Fs == nil {
		d.Fs = afero.NewMemMapFs()
	}
	if err = d.HtmlApp.WriteRef(d.Ref); err != nil {
		return
	}
	ref.Path = path.Join(ref.Path, d.BundleName())
	d.Zip.Unpacked = d.Fs
	return d.Zip.WriteRef(ref)
}

type Html struct {
	IndexHtml Blob `json:"index,omitempty" yaml:"index,omitempty"`
}

func (h *Html) ReadSrc(src Source) (err error) {
	ref := *src.Ref_()
	ref.Path = path.Join(ref.Path, "index.html")
	return h.IndexHtml.ReadSrc(&ref)
}

func (h *Html) WriteRef(ref Ref) (err error) {
	ref.Path = path.Join(ref.Path, "index.html")
	return h.IndexHtml.WriteRef(ref)
}
