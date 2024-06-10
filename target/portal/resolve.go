package portal

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/port"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"strings"
)

type Resolver[T target.Portal] struct {
	resolve target.Resolve[T]
	source  target.Source
	prefix  string
}

func NewResolver[T target.Portal](
	resolve target.Resolve[T],
	source target.Source,
) *Resolver[T] {
	return &Resolver[T]{
		resolve: resolve,
		source:  source,
		prefix:  port.PrefixStr(),
	}
}

func (f Resolver[T]) ById(id string) (t T, err error) {
	id = strings.TrimPrefix(id, f.prefix)
	id = strings.TrimPrefix(id, ".")
	for _, t = range source.List[T](f.resolve, f.source) {
		m := t.Manifest()
		if id == m.Name || id == m.Package {
			return
		}
	}
	err = errors.New("not found")
	return
}

func (f Resolver[T]) Path(path string) (p string, err error) {
	t, err := f.ById(path)
	if err != nil {
		return
	}
	p = t.Abs()
	return
}
