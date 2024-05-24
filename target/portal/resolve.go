package portal

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"log"
)

type Resolver[T target.Portal] struct {
	resolve target.Resolve[T]
	source  target.Source
}

func NewResolver[T target.Portal](
	resolve target.Resolve[T],
	source target.Source,
) *Resolver[T] {
	return &Resolver[T]{
		resolve: resolve,
		source:  source,
	}
}

func (f Resolver[T]) ById(id string) (t T, err error) {
	for _, t = range source.List[T](f.resolve, f.source) {
		log.Println(t.Abs(), t.Manifest())
		m := t.Manifest()
		if id == m.Name || id == m.Package {
			return
		}
	}
	err = errors.New("not found")
	return
}

func (f Resolver[T]) Path(id string) (p string, err error) {
	t, err := f.ById(id)
	if err != nil {
		return
	}
	p = t.Abs()
	return
}