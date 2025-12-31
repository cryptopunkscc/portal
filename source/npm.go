package source

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"

	js "github.com/cryptopunkscc/portal/core/js/embed"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
	"go.nhat.io/aferocopy/v2"
)

type NpmProject struct {
	NodeModule
	Manifest ProjectMetadata
}

func (p NpmProject) New() (src Source) {
	return &p
}

func (p NpmProject) Project() Project {
	return Project{p.Ref, p.Manifest}
}

func (p *NpmProject) ReadSrc(src Source) (err error) {
	return Readers{&p.NodeModule, &p.Manifest}.ReadSrc(src)
}

func (p *NpmProject) WriteRef(ref Ref) (err error) {
	return Writers{&p.Manifest, &p.NodeModule}.WriteRef(ref)
}

func (p NpmProject) Clean() (err error) {
	return p.Fs.RemoveAll(path.Join(p.Path, "dist"))
}

func (p NpmProject) Build() (err error) {
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

type NodeModule struct {
	Ref
	PackageJson PackageJson
}

func (p *NodeModule) ReadSrc(src Source) (err error) {
	return Readers{&p.Ref, &p.PackageJson}.ReadSrc(src)
}

func (p *NodeModule) WriteRef(ref Ref) (err error) {
	return Writers{&p.PackageJson, &p.Ref}.WriteRef(ref)
}

func (p *NodeModule) NpmInstall() (err error) {
	cmd := exec.Command("npm", "install")
	cmd.Dir = p.Path
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *NodeModule) Build() (err error) {
	if !p.PackageJson.CanBuild() {
		return plog.Errorf("missing scripts.build definition in %s/package.json", p.Path)
	}
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = p.Path
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type PackageJson struct {
	Scripts struct {
		Build string `json:"build"`
	} `json:"scripts,omitempty"`
}

func (p *PackageJson) CanBuild() bool {
	return p.Scripts.Build != ""
}

func (p *PackageJson) ReadSrc(src Source) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := afero.ReadFile(src.Ref_().Fs, path.Join(src.Ref_().Path, "package.json"))
	if err != nil {
		return
	}
	return json.Unmarshal(bytes, p)
}

func (p *PackageJson) WriteRef(ref Ref) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := json.Marshal(p)
	if err != nil {
		return
	}
	if err = afero.WriteFile(ref.Fs, path.Join(ref.Path, "package.json"), bytes, 0644); err != nil {
		return
	}
	return
}

var SkipNodeModules = &SkipDir{"node_modules"}
