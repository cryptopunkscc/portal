package source

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
)

type Blob []byte

func (m *Blob) ReadFile(files afero.Fs, name string) (err error) {
	defer plog.TraceErr(&err)
	if *m, err = afero.ReadFile(files, name); err != nil {
		return
	}
	return
}

func (m *Blob) WriteFile(dir afero.Fs, name string) (err error) {
	defer plog.TraceErr(&err)

	if err = afero.WriteFile(dir, name, *m, 0644); err != nil {
		return
	}

	return
}
