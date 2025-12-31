package source

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
)

type JsProject struct {
	NpmProject
	Js   Js
	Dist JsDist
}

func (p *JsProject) ReadSource(source Source) (err error) {
	p.Source = source
	return p.ReadFs(source.Fs)
}

func (p *JsProject) ReadFs(files afero.Fs) (err error) {
	return FSReaders{&p.NpmProject, &p.Js}.ReadFs(files)
}

func (p *JsProject) WriteFs(dir afero.Fs) (err error) {
	return FsWriters{&p.NpmProject, &p.Js}.WriteFs(dir)
}

func (p *JsProject) WriteOS(dir string) (err error) {
	return p.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type JsDist struct {
	Dist
	Js
}

func (d *JsDist) ReadSource(source Source) (err error) {
	d.Source = source
	return d.ReadFs(source.Fs)
}

func (d *JsDist) ReadFs(dir afero.Fs) (err error) {
	return FSReaders{&d.Dist, &d.Js}.ReadFs(dir)
}

func (d *JsDist) WriteFs(dir afero.Fs) (err error) {
	return FsWriters{&d.Dist, &d.Js}.WriteFs(dir)
}

func (d *JsDist) WriteOS(dir string) (err error) {
	return d.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type JsBundle struct {
	DistBundle
	Js Js
}

func (d *JsBundle) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = d.DistBundle.ReadFs(files); err != nil {
		return
	}
	if err = d.Js.ReadFs(d.ZipFs); err != nil {
		return
	}
	return
}

func (d *JsBundle) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	mapFs := afero.NewMemMapFs()
	if err = d.Js.WriteFs(mapFs); err != nil {
		return
	}
	if err = d.Dist.WriteFs(mapFs); err != nil {
		return
	}

	d.DistBundle.Fs = mapFs
	if err = d.DistBundle.WriteFs(dir); err != nil {
		return
	}
	return
}

func (d *JsBundle) WriteOS(dir string) (err error) {
	return d.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type Js struct {
	MainJs Blob `json:"main,omitempty" yaml:"main,omitempty"`
}

func (j *Js) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = j.MainJs.ReadFile(files, "main.js"); err != nil {
		return
	}
	return
}

func (j *Js) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = j.MainJs.WriteFile(dir, "main.js"); err != nil {
		return
	}
	return
}
