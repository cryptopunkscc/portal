package npm

import (
	"path"

	js "github.com/cryptopunkscc/portal/core/js/embed"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/spf13/afero"
	"go.nhat.io/aferocopy/v2"
)

type Project struct {
	NodeModule
	app.ProjectMetadata
}

func (p Project) New() (src source.Source) {
	return &p
}

func (p Project) GetMetadata() (src app.Metadata) {
	return p.Metadata
}

func (p Project) Project() app.Project {
	return app.Project{p.Ref, p.ProjectMetadata}
}

func (p Project) Clean() (err error) {
	return p.Fs.RemoveAll(path.Join(p.Path, "dist"))
}

func (p Project) Build() (err error) {
	// TODO skip if project not changed

	if err = p.NpmInstall(); err != nil {
		return
	}

	// copy portal JS lib to project node_modules
	if err = aferocopy.Copy(
		"portal",
		path.Join(p.Path, "node_modules", "portal"),
		aferocopy.Options{
			SrcFs:             afero.FromIOFS{FS: js.PortalLibFS},
			DestFs:            p.Fs,
			PermissionControl: aferocopy.AddPermission(0644),
		},
	); err != nil {
		return
	}

	// build app sources into the dist directory
	if err = p.NodeModule.Build(); err != nil {
		return
	}

	return p.Project().Build("dist")
}

func (p *Project) ReadSrc(src source.Source) (err error) {
	return source.Readers{&p.NodeModule, &p.ProjectMetadata}.ReadSrc(src)
}

func (p *Project) WriteRef(ref source.Ref) (err error) {
	return source.Writers{&p.ProjectMetadata, &p.NodeModule}.WriteRef(ref)
}
