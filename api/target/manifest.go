package target

import (
	"encoding/json"
	"io"
	"strings"
)

const ManifestFilename = "portal"

type Manifest struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Package     string `json:"package,omitempty" yaml:"package,omitempty"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
	Exec        string `json:"exec,omitempty" yaml:"exec,omitempty"`
	Schema      string `json:"schema,omitempty" yaml:"schema,omitempty"`
	OS          string `json:"os,omitempty" yaml:"os,omitempty"`
	Arch        string `json:"arch,omitempty" yaml:"arch,omitempty"`
	Hidden      bool   `json:"hidden,omitempty" yaml:"hidden,omitempty"`
	Env         Env    `json:"env,omitempty" yaml:"env,omitempty"`
}

type Env struct {
	Timeout int64 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (m Manifest) Match(id string) bool {
	return id == m.Name || strings.HasPrefix(id, m.Package)
}

func (Manifest) ObjectType() string { return "app.manifest" }

func (m Manifest) WriteTo(w io.Writer) (n int64, err error) {
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	nn, err := w.Write(b)
	if err != nil {
		return
	}
	n += int64(nn)
	return
}

func (m *Manifest) ReadFrom(r io.Reader) (n int64, err error) {
	rr := &countingReader{Reader: r}
	err = json.NewDecoder(rr).Decode(m)
	n = rr.n
	return
}

type countingReader struct {
	io.Reader
	n int64
}

func (cr *countingReader) Read(p []byte) (int, error) {
	n, err := cr.Reader.Read(p)
	cr.n += int64(n)
	return n, err
}
