package source

import (
	"errors"
	"io/fs"
	"strings"
)

type Source interface {
	Reader
	Ref_() *Ref
}

type Constructor interface {
	Source
	New() Source
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

var SkipNodeModules = &SkipDir{"node_modules"}
