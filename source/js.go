package source

import (
	"path"

	"github.com/spf13/afero"
)

type JsProject struct {
	NpmProject
	Js  Js
	App JsApp
}

func (p JsProject) New() (src Source) {
	return &p
}

func (p *JsProject) ReadSrc(src Source) (err error) {
	return Readers{&p.Js, &p.NpmProject}.ReadSrc(src)
}

func (p *JsProject) WriteRef(ref Ref) (err error) {
	return Writers{&p.NpmProject, &p.Js}.WriteRef(ref)
}

type JsApp struct {
	App
	Js
}

func (a JsApp) New() (src Source) {
	return &a
}

func (a *JsApp) ReadSrc(src Source) (err error) {
	return Readers{&a.App, &a.Js}.ReadSrc(src)
}

func (a *JsApp) WriteRef(ref Ref) (err error) {
	return Writers{&a.App, &a.Js}.WriteRef(ref)
}

type JsBundle struct {
	JsApp
	Zip
}

func (b JsBundle) New() (src Source) {
	return &b
}

func (b *JsBundle) ReadSrc(src Source) (err error) {
	if err = b.Zip.ReadSrc(src); err != nil {
		return
	}
	if err = b.JsApp.ReadSrc(&Ref{Fs: b.Zip.Unpacked}); err != nil {
		return
	}
	return
}

func (b *JsBundle) WriteRef(ref Ref) (err error) {
	if b.Fs == nil {
		b.Fs = afero.NewMemMapFs()
	}
	if err = b.JsApp.WriteRef(b.Ref); err != nil {
		return
	}
	ref.Path = path.Join(ref.Path, b.BundleName())
	b.Zip.Unpacked = b.Fs
	return b.Zip.WriteRef(ref)
}

type Js struct {
	MainJs Blob `json:"main,omitempty" yaml:"main,omitempty"`
}

func (j *Js) ReadSrc(src Source) (err error) {
	ref := *src.Ref_()
	ref.Path = path.Join(ref.Path, "main.js")
	return j.MainJs.ReadSrc(&ref)
}

func (j *Js) WriteRef(ref Ref) (err error) {
	ref.Path = path.Join(ref.Path, "main.js")
	return j.MainJs.WriteRef(ref)
}
