package source

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
)

type Blob []byte

func (m *Blob) ReadSrc(src Source) (err error) {
	defer plog.TraceErr(&err)
	if *m, err = afero.ReadFile(src.Ref_().Fs, src.Ref_().Path); err != nil {
		return
	}
	return
}

func (m *Blob) WriteRef(ref Ref) (err error) {
	defer plog.TraceErr(&err)
	if err = afero.WriteFile(ref.Fs, ref.Path, *m, 0644); err != nil {
		return
	}
	return
}
