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

type Refs []Ref

func (r Refs) Collect(constructors ...Constructor) (out []Source) {
	for _, ref := range r {
		if out = ref.Collect(constructors...); len(out) > 0 {
			return
		}
	}
	return
}

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

type MapFunc func(Source) (Source, error)

func (r Ref) Collect(constructors ...Constructor) (out List[Source]) {
	if r.Fs == nil {
		return
	}
	dir := r.Path
	_ = afero.Walk(r.Fs, r.Path, func(p string, _ fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		r.Path = path2.Join(dir, p)
		source, err := r.Resolve(constructors...)
		if isSkip(err) {
			return err
		}
		out = append(out, source)
		return nil
	})
	return
}

func (r Ref) Resolve(constructors ...Constructor) (out Source, err error) {
	if len(constructors) == 0 {
		return nil, fs.ErrNotExist
	}
	for _, constructor := range constructors {
		out = constructor.New()
		if err = out.ReadSrc(&r); err == nil {
			return
		}
		if isSkip(err) {
			return nil, err
		}
	}
	return nil, fs.ErrInvalid
}

func CollectIt[T interface {
	Source
	Constructor
}](src Source, constructor T) (out List[T]) {
	return CollectT[T](src, constructor)
}

func Collect(src Source, constructors ...Constructor) (out List[Source]) {
	return Collect(src, constructors...)
}

func CollectT[T Source](src Source, constructors ...Constructor) (out List[T]) {
	if src == nil {
		return
	}

	constructors = append([]Constructor{SkipNodeModules}, constructors...)

	if src != src.Ref_() {
		if o, err := ResolveT[T](src, constructors...); err == nil {
			out = append(out, o)
		}
		return
	}

	ref := *src.Ref_()
	if ref.Fs == nil {
		return
	}
	_ = afero.Walk(ref.Fs, ref.Path, func(p string, _ fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ref.Path = p
		source, err := ResolveT[T](&ref, constructors...)
		if isSkip(err) {
			return err
		}
		if err == nil {
			out = append(out, source)
		}
		return nil
	})
	return
}

func ResolveIt[T interface {
	Source
	Constructor
}](src Source, constructor T) (out T, err error) {
	return ResolveT[T](src, constructor)
}

func Resolve(src Source, constructors ...Constructor) (out Source, err error) {
	return ResolveT[Source](src, constructors...)
}

func ResolveT[T Source](src Source, constructors ...Constructor) (out T, err error) {
	if out, err = Cast[T](src); err == nil {
		return
	}
	if len(constructors) == 0 {
		return out, fs.ErrInvalid
	}
	for _, constructor := range constructors {
		if out, err = Cast[T](constructor.New()); err == nil {
			if err = out.ReadSrc(src); err == nil {
				return
			}
		}
		if isSkip(err) {
			return
		}
	}
	return out, fs.ErrInvalid
}

func Cast[T Source](src Source) (out T, err error) {
	out, ok := src.(T)
	if !ok {
		return out, fs.ErrInvalid
	}
	return
}

func isSkip(err error) bool {
	return errors.Is(err, fs.SkipDir) || errors.Is(err, fs.SkipAll)
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

type Reader interface {
	ReadSrc(src Source) (err error)
}

type Readers []Reader

func (r Readers) ReadSrc(src Source) (err error) {
	for _, reader := range r {
		if err = reader.ReadSrc(src); err != nil {
			return
		}
	}
	return
}

type Writer interface {
	WriteRef(ref Ref) (err error)
}

type Writers []Writer

func (w Writers) WriteRef(ref Ref) (err error) {
	for _, writer := range w {
		if err = writer.WriteRef(ref); err != nil {
			return
		}
	}
	return
}

type List[T any] []T

func (l List[T]) Filter(f func(T) bool) (out List[T]) {
	for _, t := range l {
		if f(t) {
			out = append(out, t)
		}
	}
	return
}

type Source interface {
	Reader

	Ref_() *Ref
}

type Constructor interface {
	New() Source
}

type Builder[C Constructor] struct {
	Constructors []C
	Source       []C
}

func (b *Builder[C]) ReadSrc(src Source) (err error) {
	ref := *src.Ref_()
	dir := src.Ref_().Path
	return afero.Walk(ref.Fs, src.Ref_().Path, func(p string, _ fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		src.Ref_().Path = path2.Join(dir, p)
		for _, c := range b.Constructors {
			source := c.New()
			if err = source.ReadSrc(ref.New()); isSkip(err) {
				return err
			}
			b.Source = append(b.Source, source.(C))
		}
		return nil
	})
}

type Filter struct {
	Func func(Ref) error
	ref  *Ref
}

func (f Filter) New() Source { return &f }
func (f *Filter) Ref_() *Ref { return f.ref }
func (f *Filter) ReadSrc(src Source) (err error) {
	if err = f.Func(*f.ref); err != nil {
		f.ref = src.Ref_()
	}
	return
}

type SkipDir struct{ Name string }

func (s SkipDir) New() Source { return &s }
func (s SkipDir) Ref_() *Ref  { return nil }
func (s SkipDir) ReadSrc(src Source) (err error) {
	if strings.HasSuffix(src.Ref_().Path, s.Name) {
		return fs.SkipDir
	}
	return errors.New(s.Name)
}
