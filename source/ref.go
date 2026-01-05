package source

import (
	"errors"
	"io/fs"
	"os"
	path2 "path"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
	"go.nhat.io/aferocopy/v2"
)

type Ref struct {
	Fs   afero.Fs
	Path string
	Func any
}

func Abs(path ...string) string {
	src := path2.Join(path...)
	if path2.IsAbs(src) {
		return src
	}
	base, err := os.Getwd()
	if err != nil {
		return src
	}
	base = filepath.ToSlash(base)
	return path2.Join(base, src)
}

func FSRef(fs fs.FS, path ...string) *Ref {
	return &Ref{Fs: afero.FromIOFS{FS: fs}, Path: path2.Join(path...)}
}

func OSRef(path ...string) *Ref {
	return &Ref{Fs: afero.NewOsFs(), Path: Abs(path...)}
}

func (r Ref) Sub(path string) *Ref {
	r.Path = path2.Join(r.Path, path)
	return &r
}

func (r Ref) New() Source {
	return &r
}

func (r Ref) String() string {
	return r.Path
}

func (r Ref) FS() fs.FS {
	return afero.IOFS{Fs: r.Fs}
}

func (r *Ref) Ref_() *Ref {
	return r
}

func (r *Ref) Checkout(path string) (err error) {
	path = path2.Join(r.Path, path)
	_, err = r.Fs.Stat(path)
	if err != nil {
		return
	}
	r.Path = path
	return
}

func (r *Ref) ReadOS(path string) (err error) {
	path = Abs(path)
	return r.ReadSrc(&Ref{Fs: afero.NewOsFs(), Path: path})
}

func (r *Ref) ReadFs(fs afero.Fs) (err error) {
	return r.ReadSrc(&Ref{Fs: fs})
}

func (r *Ref) ReadSrc(src Source) (err error) {
	defer plog.TraceErr(&err)
	ref := *src.Ref_()
	if ref.Fs == nil {
		return errors.New("Ref.ReadSrc: ref.Fs cannot be nil")
	}
	p := "."
	if src.Ref_().Path != "" {
		p = src.Ref_().Path
	}
	if _, err = ref.Fs.Stat(p); err != nil {
		return
	}
	*r = ref
	return
}

func (r *Ref) WriteOS(dir string) (err error) {
	return r.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

func (r *Ref) WriteFs(fs afero.Fs) (err error) {
	return r.WriteRef(Ref{Fs: fs})
}

func (r *Ref) WriteRef(ref Ref) (err error) {
	if r.Fs == nil || r.Fs == ref.Fs {
		return
	}
	defer plog.TraceErr(&err)
	srcPath := "."
	if len(r.Path) > 0 {
		srcPath = r.Path
	}
	err = aferocopy.Copy(srcPath, r.Ref_().Path, aferocopy.Options{
		SrcFs: r.Fs, DestFs: ref.Fs,
		Skip: func(srcFs afero.Fs, src string) (bool, error) {
			return strings.HasPrefix(src, "/"), nil
		},
	})
	if err == nil {
		*r = ref
	}
	return
}
