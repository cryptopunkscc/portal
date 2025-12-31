package source

import (
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
)

type HtmlProject struct {
	Project
	Html        Html
	PackageJson PackageJson

	Dist HtmlDist
}

func (p *HtmlProject) ReadFs(files afero.Fs) (err error) {
	return FSReaders{&p.Project, &p.Html, &p.PackageJson}.ReadFs(files)
}

func (p *HtmlProject) WriteFs(dir afero.Fs) (err error) {
	return FsWriters{&p.Project, &p.Html, &p.PackageJson}.WriteFs(dir)
}

func (p *HtmlProject) WriteOS(dir string) (err error) {
	return p.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type HtmlDist struct {
	Dist
	Html
}

func (d *HtmlDist) ReadFs(files afero.Fs) (err error) {
	return FSReaders{&d.Dist, &d.Html}.ReadFs(files)
}

func (d *HtmlDist) WriteFs(dir afero.Fs) (err error) {
	return FsWriters{&d.Dist, &d.Html}.WriteFs(dir)
}

func (d *HtmlDist) WriteOS(dir string) (err error) {
	return d.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type HtmlBundle struct {
	DistBundle
	Html Html
}

func (d *HtmlBundle) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = d.DistBundle.ReadFs(files); err != nil {
		return
	}
	if err = d.Html.ReadFs(d.ZipFs); err != nil {
		return
	}
	return
}

func (d *HtmlBundle) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = d.Html.WriteFs(dir); err != nil {
		return
	}
	if err = d.DistBundle.WriteFs(dir); err != nil {
		return
	}
	return
}

func (d *HtmlBundle) WriteOS(dir string) (err error) {
	return d.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type Html struct {
	IndexHtml Blob `json:"index,omitempty" yaml:"index,omitempty"`
}

func (h *Html) ReadFS(files fs.FS) (err error) {
	return h.ReadFs(afero.FromIOFS{FS: files})
}

func (h *Html) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = h.IndexHtml.ReadFile(files, "index.html"); err != nil {
		return
	}
	return
}

func (h *Html) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = h.IndexHtml.WriteFile(dir, "index.html"); err != nil {
		return
	}
	return
}
