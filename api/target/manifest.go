package target

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"io"
	"strings"
)

const ManifestFilename = "portal"

type Manifest struct {
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
	Title        string `json:"title,omitempty" yaml:"title,omitempty"`
	Description  string `json:"description,omitempty" yaml:"description,omitempty"`
	Package      string `json:"package,omitempty" yaml:"package,omitempty"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
	Icon         string `json:"icon,omitempty" yaml:"icon,omitempty"`
	Schema       string `json:"schema,omitempty" yaml:"schema,omitempty"`
	Hidden       bool   `json:"hidden,omitempty" yaml:"hidden,omitempty"`
	Env          Env    `json:"env,omitempty" yaml:"env,omitempty"`
	Distribution `json:",omitempty,inline" yaml:",omitempty,inline"`
}

type Env struct {
	Timeout int64 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (m Manifest) Match(id string) bool {
	return id == m.Name || strings.HasPrefix(id, m.Package)
}

func (Manifest) ObjectType() string { return "app.manifest" }

func (m Manifest) WriteTo(w io.Writer) (n int64, err error) { return WriteTo(w, json.Marshal, m) }

func (m *Manifest) ReadFrom(r io.Reader) (n int64, err error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return
	}
	if err = all.Unmarshalers.Unmarshal(b, m); err != nil {
		return
	}
	n = int64(len(b))
	return
}

func WriteTo(w io.Writer, marshal stream.Marshal, a any) (n int64, err error) {
	b, err := marshal(a)
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

func ReadJsonFrom(r io.Reader, a any) (n int64, err error) {
	rr := &countingReader{Reader: r}
	err = json.NewDecoder(rr).Decode(a)
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

type Distribution struct {
	Exec string `json:"exec,omitempty" yaml:"exec,omitempty"`
	OS   string `json:"os,omitempty" yaml:"os,omitempty"`
	Arch string `json:"arch,omitempty" yaml:"arch,omitempty"`
}
