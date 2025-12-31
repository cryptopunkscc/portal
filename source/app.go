package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type App struct {
	Ref
	Metadata Metadata
}

func (a *App) BundleName() string {
	t := a.Metadata.Target
	platform := ""
	if len(t.OS) > 0 {
		platform += "_" + t.OS
		if len(t.Arch) > 0 {
			platform += "_" + t.Arch
		}
	}

	m := a.Metadata
	version := fmt.Sprintf("%d.%d.%d", m.Version, m.Api.Version, m.Release.Version)

	return fmt.Sprintf("%s_%s%s.portal", m.Package, version, platform)
}

func (a *App) ReadSrc(src Source) (err error) {
	return Readers{&a.Ref, &a.Metadata}.ReadSrc(src)
}

func (a *App) WriteRef(ref Ref) (err error) {
	return Writers{&a.Ref, &a.Metadata}.WriteRef(ref)
}

func (a App) Bundle() *AppBundle {
	return &AppBundle{App: a, Zip: Zip{
		Unpacked: afero.NewBasePathFs(a.Fs, a.Path)}}
}

type AppBundle struct {
	App
	Zip
}

func init() {
	_ = astral.DefaultBlueprints.Add(&AppBundle{})
}

func (b *AppBundle) ReadSrc(src Source) (err error) {
	if err = b.Zip.ReadSrc(src); err != nil {
		return
	}
	if err = b.App.ReadFs(b.Unpacked); err != nil {
		return
	}
	return
}

func (b *AppBundle) WriteRef(ref Ref) (err error) {
	if b.App.Fs == nil {
		b.App.Fs = afero.NewMemMapFs()
	}
	if err = b.App.WriteRef(b.App.Ref); err != nil {
		return
	}

	b.Zip.Unpacked = b.App.Fs
	if len(b.App.Path) > 0 {
		b.Zip.Unpacked = afero.NewBasePathFs(b.Fs, b.Path)
	}

	ref.Path = path.Join(ref.Path, b.BundleName())
	return b.Zip.WriteRef(ref)
}

func (b *AppBundle) ObjectType() string { return "app.bundle" }

func (b *AppBundle) WriteTo(w io.Writer) (n int64, err error) {
	defer plog.TraceErr(&err)
	i, err := w.Write(b.Blob)
	return int64(i), err
}

func (b *AppBundle) ReadFrom(r io.Reader) (n int64, err error) {
	defer plog.TraceErr(&err)
	blob, err := io.ReadAll(r)
	if err != nil {
		return
	}
	return b.Zip.ReadFrom(bytes.NewReader(blob))
}

func (b AppBundle) Publish(objects *astrald.ObjectsClient) (err error) {
	if err = b.Zip.Publish(objects); err != nil {
		return
	}

	release := ReleaseMetadata{
		BundleID: b.Zip.ObjectID,
		Release:  b.Metadata.Release,
		Target:   b.Metadata.Target,
	}
	if release.ManifestID, err = astral.ResolveObjectID(&b.Metadata.Manifest); err != nil {
		return
	}

	return objects.Push(nil, &release, &b.Metadata.Manifest)
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

func init() {
	_ = astral.DefaultBlueprints.Add(&Manifest{})
}

func (m *Manifest) Match(id string) bool {
	return id == m.Name || strings.HasPrefix(id, m.Package)
}

func (m *Manifest) ObjectType() string { return "app.manifest" }

func (m *Manifest) ReadSrc(src Source) (err error) {
	b, err := afero.ReadFile(src.Ref_().Fs, path.Join(src.Ref_().Path, "portal.json"))
	if err == nil {
		return json.Unmarshal(b, m)
	}
	b, err = afero.ReadFile(src.Ref_().Fs, path.Join(src.Ref_().Path, "portal.yaml"))
	if err == nil {
		return yaml.Unmarshal(b, m)
	}
	b, err = afero.ReadFile(src.Ref_().Fs, path.Join(src.Ref_().Path, "portal.yml"))
	if err == nil {
		return yaml.Unmarshal(b, m)
	}
	return
}

func (m *Manifest) WriteRef(ref Ref) (err error) {
	defer plog.TraceErr(&err)
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	return afero.WriteFile(ref.Fs, path.Join(ref.Path, "portal.json"), b, 0644)
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

func (m *Manifest) ReadFrom(r io.Reader) (n int64, err error) {
	defer plog.TraceErr(&err)
	b, err := io.ReadAll(r)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, m); err == nil {
		return int64(len(b)), nil
	}
	if err2 := yaml.Unmarshal(b, m); err2 == nil {
		return int64(len(b)), nil
	}
	return
}

func (m *Manifest) WriteTo(w io.Writer) (n int64, err error) {
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	nn, err := w.Write(b)
	if err != nil {
		return
	}
	n = int64(nn)
	return
}
