package source

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
