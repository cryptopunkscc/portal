package app

import (
	"encoding/json"
	"io"
	"path"
	"runtime"
	"strings"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/source"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type App interface {
	source.Source
	GetPath() string
	GetDist() Dist
	GetMetadata() Metadata
}

type Metadata struct {
	Manifest `json:",inline" yaml:",inline"`
	Api      Api     `json:"api,omitempty" yaml:"api,omitempty"`
	Config   Config  `json:"config,omitempty" yaml:"config,omitempty"`
	Target   Target  `json:"target,omitempty" yaml:"target,omitempty"`
	Release  Release `json:"release,omitempty" yaml:"release,omitempty"`
}

func (m *Metadata) ReadSrc(src source.Source) (err error) {
	return metadataReadSrc(m, "portal", src)
}

func (m *Metadata) WriteRef(ref source.Ref) (err error) {
	return metadataWriteRef(m, "portal", ref)
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

func (m *Manifest) ReadSrc(src source.Source) (err error) {
	return metadataReadSrc(m, "portal", src)
}

func (m *Manifest) WriteRef(ref source.Ref) (err error) {
	return metadataWriteRef(m, "portal", ref)
}

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

func (r Target) Match() bool { return r.OS == runtime.GOOS && r.Arch == runtime.GOARCH }

type Release struct {
	Version int `json:"version,omitempty" yaml:"version,omitempty"`
}

func metadataWriteRef(meta any, name string, ref source.Ref) (err error) {
	defer plog.TraceErr(&err)
	b, err := json.Marshal(meta)
	if err != nil {
		return
	}
	return afero.WriteFile(ref.Fs, path.Join(ref.Path, name+".json"), b, 0644)
}

func metadataReadSrc(meta any, name string, src source.Source) (err error) {
	defer plog.TraceErr(&err)
	ref := *src.Ref_()
	b, err := afero.ReadFile(ref.Fs, path.Join(ref.Path, name+".json"))
	if err == nil {
		return json.Unmarshal(b, meta)
	}
	b, err = afero.ReadFile(ref.Fs, path.Join(ref.Path, name+".yaml"))
	if err == nil {
		return yaml.Unmarshal(b, meta)
	}
	b, err = afero.ReadFile(ref.Fs, path.Join(ref.Path, name+".yml"))
	if err == nil {
		return yaml.Unmarshal(b, meta)
	}
	return
}
