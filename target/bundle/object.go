package bundle

import (
	"io"
	"testing/fstest"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/source"
)

type Object[T any] struct {
	target.Bundle_
	target.AppBundle[T]
	target.Resolve[target.AppBundle[T]]
}

var _ astral.Object = &Object[any]{}

func (o *Object[T]) ObjectType() string { return "app.bundle" }

func (o *Object[T]) WriteTo(w io.Writer) (n int64, err error) {
	if o.Bundle_ == nil {
		o.Bundle_ = o.AppBundle
	}
	f, err := o.Bundle_.Package().File()
	if err != nil {
		return
	}
	return io.Copy(w, f)
}

func (o *Object[T]) ReadFrom(r io.Reader) (n int64, err error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return
	}
	n = int64(len(b))

	f := fstest.MapFS{
		o.ObjectType(): &fstest.MapFile{
			Data: b,
		},
	}
	s, err := source.Embed(f).Sub(o.ObjectType())
	if err != nil {
		return
	}
	o.AppBundle, err = o.Resolve(s)
	return
}
