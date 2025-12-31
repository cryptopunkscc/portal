package source

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
)

type FsReader interface {
	ReadFs(files afero.Fs) (err error)
}

type FSReaders []FsReader

func (f FSReaders) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	for _, reader := range f {
		if err = reader.ReadFs(files); err != nil {
			return
		}
	}
	return
}

type FsWriter interface {
	WriteFs(dir afero.Fs) (err error)
}

type FsWriters []FsWriter

func (w FsWriters) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	for _, writer := range w {
		if err = writer.WriteFs(dir); err != nil {
			return
		}
	}
	return
}

type FsReadWriter interface {
	FsReader
	FsWriter
}
