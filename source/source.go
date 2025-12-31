package source

import (
	"archive/zip"
	"bytes"
	"errors"
	"strings"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
	"go.nhat.io/aferocopy/v2"
)

type Source struct {
	Fs   afero.Fs
	Name string
	Dir  string
}

func (s *Source) ReadFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if dir != nil {
		s.Fs = dir
	}
	if s.Fs == nil {
		return errors.New("the Source.FS is not initialized")
	}
	return nil
}

func (s *Source) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if s.Fs != nil && s.Fs != dir {
		if s.Name == "" {
			s.Name = "."
		}
		if err = aferocopy.Copy(s.Name, "", aferocopy.Options{
			SrcFs:  s.Fs,
			DestFs: dir,
			Skip: func(srcFs afero.Fs, src string) (bool, error) {
				return strings.HasPrefix(src, "/"), nil
			},
		}); err != nil {
			return
		}
	}
	s.Fs = dir
	return
}

func (s *Source) WriteOS(dir string) (err error) {
	return s.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

func (s *Source) WriteZipFs(out afero.Fs) (err error) {
	defer plog.TraceErr(&err)

	buffer := bytes.Buffer{}
	zipWriter := zip.NewWriter(&buffer)

	if err = zipWriter.AddFS(afero.IOFS{Fs: s.Fs}); err != nil {
		return
	}
	if err = zipWriter.Close(); err != nil {
		return
	}

	if out == nil {
		out = afero.NewMemMapFs()
	} else {
		s.Fs = out
	}

	if err = afero.WriteFile(out, s.Name, buffer.Bytes(), 0644); err != nil {
		return
	}

	return
}
