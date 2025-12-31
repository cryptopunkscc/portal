package source

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/spf13/afero"
)

type Dist struct {
	Source
	Metadata Metadata
}

func (a *Dist) ReadFs(dir afero.Fs) (err error) {
	return FSReaders{&a.Source, &a.Metadata}.ReadFs(dir)
}

func (a *Dist) WriteFs(dir afero.Fs) (err error) {
	return FsWriters{&a.Metadata, &a.Source}.WriteFs(dir)
}

type DistBundle struct {
	ZipBundle
	Dist Dist
}

func (b *DistBundle) ReadSource(source Source) (err error) {
	b.Source = source
	return b.ReadFs(source.Fs)
}

func (b *DistBundle) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = b.ZipBundle.ReadFs(files); err != nil {
		return
	}
	if err = b.Dist.ReadFs(b.ZipFs); err != nil {
		return
	}
	return
}

func (b *DistBundle) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = b.Dist.WriteFs(b.Dist.Fs); err != nil {
		return
	}

	t := b.Dist.Metadata.Target
	platform := ""
	if len(t.OS) > 0 {
		platform += "_" + t.OS
		if len(t.Arch) > 0 {
			platform += "_" + t.Arch
		}
	}

	m := b.Dist.Metadata
	version := fmt.Sprintf("%d.%d.%d", m.Version, m.Api.Version, m.Release.Version)

	b.ZipBundle.Name = fmt.Sprintf("%s_%s%s.portal", b.Dist.Metadata.Package, version, platform)
	b.ZipBundle.Fs = b.Dist.Fs
	return b.ZipBundle.WriteZipFs(dir)
}

func (b *DistBundle) WriteOS(dir string) (err error) {
	return b.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type Metadata struct {
	Manifest `json:",inline" yaml:",inline"`
	Api      Api     `json:"api,omitempty" yaml:"api,omitempty"`
	Config   Config  `json:"config,omitempty" yaml:"config,omitempty"`
	Target   Target  `json:"target,omitempty" yaml:"target,omitempty"`
	Release  Release `json:"release,omitempty" yaml:"release,omitempty"`
}

type Manifest struct {
	// Name it a sort name of the application. No space allowed. Lowercased.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Title is a full name of the application used to display in GUI.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// Description of the Application.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Package name in reverse domain style.
	Package string `json:"package,omitempty" yaml:"package,omitempty"`
	// Version number of the App. 0 - pre releases. 1 - first release. After first release any change should increment Version number.
	Version int `json:"version,omitempty" yaml:"version,omitempty"`
	// Icon path relative to the App.
	Icon string `json:"icon,omitempty" yaml:"icon,omitempty"`
	// Runtime of the application [js, html, exec].
	Runtime string `json:"schema,omitempty" yaml:"schema,omitempty"`
	// Type of the application [gui, cli, api].
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

func (m *Manifest) ReadFs(files afero.Fs) (err error) {
	bytes, err := afero.ReadFile(files, "portal.json")
	if err == nil {
		return json.Unmarshal(bytes, m)
	}
	bytes, err = afero.ReadFile(files, "portal.yaml")
	if err == nil {
		return json.Unmarshal(bytes, m)
	}
	return
}

func (m *Manifest) WriteFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := json.Marshal(m)
	if err != nil {
		return
	}
	return afero.WriteFile(files, "portal.json", bytes, 0644)
}

type Api struct {
	Version     int `json:"version,omitempty" yaml:"version,omitempty"`
	cmd.Handler `json:",omitempty" yaml:",omitempty,inline"`
}

type Config struct {
	Timeout int64 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Hidden  bool  `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

type Target struct {
	Exec string `json:"exec,omitempty" yaml:"exec,omitempty"`
	OS   string `json:"os,omitempty" yaml:"os,omitempty"`
	Arch string `json:"arch,omitempty" yaml:"arch,omitempty"`
}

type Release struct {
	Version int `json:"version,omitempty" yaml:"version,omitempty"`
}

func (r Target) Match() bool { return r.OS == runtime.GOOS && r.Arch == runtime.GOARCH }
