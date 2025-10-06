package manifest

import (
	"encoding/json"
	"io"
	"io/fs"
	"strings"

	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

const AppFilename = "portal"

type App struct {
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

func (a *App) UnmarshalFrom(bytes []byte) error { return all.Unmarshalers.Unmarshal(bytes, a) }
func (a *App) LoadFrom(fs fs.FS, files ...string) error {
	if len(files) == 0 {
		files = []string{AppFilename, DevFilename}
	}
	return all.Unmarshalers.Load(a, fs, files...)
}

func (a App) Match(id string) bool {
	return id == a.Name || strings.HasPrefix(id, a.Package)
}

func (App) ObjectType() string { return "app.manifest" }

func (a *App) ReadFrom(r io.Reader) (n int64, err error) {
	defer plog.TraceErr(&err)
	b, err := io.ReadAll(r)
	if err != nil {
		return
	}
	if err = all.Unmarshalers.Unmarshal(b, a); err != nil {
		return
	}
	n = int64(len(b))
	return
}

func (a App) WriteTo(w io.Writer) (n int64, err error) {
	b, err := json.Marshal(a)
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
